package evaluator

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Wa4h1h/memdb/internal/utils"
	"go.uber.org/zap"

	"github.com/Wa4h1h/memdb/internal/store"
)

type Evaluation struct {
	command string
	key     string
	value   string
	ttl     int64
}

type Evaluator struct {
	store.Store
	l *zap.SugaredLogger
}

func NewEvaluator(s store.Store, l *zap.SugaredLogger) *Evaluator {
	return &Evaluator{s, l}
}

func (e *Evaluator) Evaluate(cmd string) (string, error) {
	lexer := NewLexer(strings.TrimSuffix(cmd, "\n"))

	var tokens []*Token

	for tk := lexer.NextToken(); tk.Literal != EOF; tk = lexer.NextToken() {
		tokens = append(tokens, tk)
	}

	evaluation, err := e.evaluateTokens(tokens)
	if err != nil {
		return "", err
	}

	switch evaluation.command {
	case SET:
		return e.set(evaluation.key, evaluation.value, evaluation.ttl)
	case GET:
		return e.get(evaluation.key)
	case REMOVE:
		return e.remove(evaluation.key)
	case INC:
		return e.inc(evaluation.key, evaluation.value)
	case DEC:
		return e.dec(evaluation.key, evaluation.value)
	case NEGATE:
		return e.negate(evaluation.key)
	case LAPPEND:
		return e.lappend(evaluation.key, evaluation.value)
	case LREMOVE:
		return e.lremove(evaluation.key, evaluation.value)
	}

	return "", utils.ErrUnknownCommand
}

func (e *Evaluator) evaluateTokens(tokens []*Token) (*Evaluation, error) {
	if len(tokens) > 5 || tokens[0].Type != COMMAND {
		return nil, utils.ErrUnknownCommand
	}

	if len(tokens) < 2 {
		return nil, utils.ErrMissingOptions
	}

	var evaluation Evaluation

	for i, token := range tokens {
		switch token.Type {
		case COMMAND:
			evaluation.command = token.Literal
		case TTL:
			val, err := strconv.ParseInt(token.Literal, 10, 64)
			if err != nil {
				e.l.Error(fmt.Sprintf("failed to pares str to int: %s", err.Error()))

				return nil, utils.ErrParsingTTL
			}

			evaluation.ttl = val
		default:
			if i == 1 {
				evaluation.key = token.Literal
			}
			if i == 2 {
				evaluation.value = token.Literal
			}
		}
	}

	return &evaluation, nil
}

func (e *Evaluator) set(key string, value string, ttl int64) (string, error) {
	item := &store.Item{
		Value:     value,
		TTL:       ttl,
		CreatedAt: time.Now(),
	}

	e.AddItem(key, item)

	_, err := e.GetItem(key)
	if err != nil {
		return "", err
	}

	return "Ok\n", nil
}

func (e *Evaluator) get(key string) (string, error) {
	item, err := e.GetItem(key)
	if err != nil {
		return "", err
	}

	return FormatString(item.Value), nil
}

func (e *Evaluator) remove(key string) (string, error) {
	val, err := e.DeleteItem(key)
	if err != nil {
		return "", err
	}

	return FormatString(val), nil
}

func (e *Evaluator) updateIntValue(key string, value string, command string) (string, error) {
	var val int64 = 1

	if len(value) > 0 {
		parsedValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			e.l.Error(fmt.Sprintf("failed to pares str to int: %s", err.Error()))

			return "", utils.ErrParsingToInt
		}

		val = parsedValue
	}

	item, err := e.GetItem(key)
	if err != nil {
		return "", err
	}

	itemValue, err := strconv.ParseInt(item.Value, 10, 64)
	if err != nil {
		e.l.Error(fmt.Sprintf("failed to pares str to int: %s", err.Error()))

		return "", utils.ErrParsingToInt
	}

	switch command {
	case INC:
		itemValue += val
	case DEC:
		itemValue -= val
	}

	item.Value = strconv.Itoa(int(itemValue))

	return e.set(key, item.Value, item.TTL)
}

func (e *Evaluator) inc(key string, value string) (string, error) {
	return e.updateIntValue(key, value, INC)
}

func (e *Evaluator) dec(key string, value string) (string, error) {
	return e.updateIntValue(key, value, DEC)
}

func (e *Evaluator) negate(key string) (string, error) {
	item, err := e.GetItem(key)
	if err != nil {
		return "", err
	}

	val, err := strconv.ParseBool(item.Value)
	if err != nil {
		e.l.Error(fmt.Sprintf("failed to pares str to bool: %s", err.Error()))

		return "", utils.ErrParsingToBool
	}

	item.Value = strconv.FormatBool(!val)

	return e.set(key, item.Value, item.TTL)
}

func (e *Evaluator) updateList(key string, value string, command string) (string, error) {
	item, err := e.GetItem(key)
	if err != nil {
		return "", err
	}

	slice, err := utils.ParseStringToSlice(item.Value)
	if err != nil {
		return "", utils.ErrMalformedSlice
	}

	switch command {
	case LAPPEND:
		{
			slice = append(slice, value)
			item.Value = utils.ParseSliceToString(slice)
		}
	case LREMOVE:
		{
			itemIndex := slices.Index(slice, value)

			if itemIndex == -1 {
				return "", utils.ErrElementNotinList
			}

			slice[itemIndex] = slice[len(slice)-1]
			item.Value = utils.ParseSliceToString(slice[:len(slice)-1])
		}

	}

	return e.set(key, item.Value, item.TTL)
}

func (e *Evaluator) lappend(key string, value string) (string, error) {
	return e.updateList(key, value, LAPPEND)
}

func (e *Evaluator) lremove(key string, value string) (string, error) {
	return e.updateList(key, value, LREMOVE)
}
