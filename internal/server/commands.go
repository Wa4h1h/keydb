package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Wa4h1h/memdb/internal/store"
	"github.com/Wa4h1h/memdb/pkg"
)

type Evaluation struct {
	command string
	key     string
	value   string
	ttl     int64
}

func (s *Server) evaluateTokens(tokens []*Token) (*Evaluation, error) {
	if len(tokens) > 5 || tokens[0].Type != COMMAND {
		return nil, pkg.ErrUnknownCommand
	}

	if len(tokens) < 2 {
		return nil, pkg.ErrMissingOptions
	}

	var evaluation Evaluation

	for i, token := range tokens {
		switch token.Type {
		case COMMAND:
			evaluation.command = token.Literal
		case TTL:
			val, err := strconv.ParseInt(token.Literal, 10, 64)
			if err != nil {
				s.Logger.Error(fmt.Sprintf("failed to pares str to int: %s", err.Error()))
				return nil, pkg.ErrParsingTTL
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

func (s *Server) execute(cmd string) (string, error) {
	lexer := NewLexer(strings.TrimSuffix(cmd, "\n"))
	var tokens []*Token
	for tk := lexer.NextToken(); tk.Literal != EOF; tk = lexer.NextToken() {
		tokens = append(tokens, tk)
	}

	evaluation, err := s.evaluateTokens(tokens)
	if err != nil {
		return "", err
	}

	switch evaluation.command {
	case SET:
		return s.set(evaluation.key, evaluation.value, evaluation.ttl)
	case GET:
		return s.get(evaluation.key)
	case REMOVE:
		return s.remove(evaluation.key)
	case INC:
		return s.inc(evaluation.key, evaluation.value)
	case DEC:
		return s.dec(evaluation.key, evaluation.value)
	case NEGATE:
		return s.negate(evaluation.key)
	case LAPPEND:
		return s.lappend(evaluation.key, evaluation.value)
	case LREMOVE:
		return s.lremove(evaluation.key, evaluation.value)
	}

	return "", pkg.ErrUnknownCommand
}

func (s *Server) set(key string, value string, ttl int64) (string, error) {
	item := &store.Item{
		Value: value,
		TTL:   ttl,
	}

	s.store.AddItem(key, item)

	_, err := s.store.GetItem(key)
	if err != nil {
		return "", err
	}

	return "Ok\n", nil
}

func (s *Server) get(key string) (string, error) {
	item, err := s.store.GetItem(key)
	if err != nil {
		return "", err
	}

	return formatString(item.Value), nil
}

func (s *Server) remove(key string) (string, error) {
	val, err := s.store.DeleteItem(key)
	if err != nil {
		return "", err
	}

	return formatString(val), nil
}

func (s *Server) updateIntValue(key string, value string, command string) (string, error) {
	var val int64 = 1
	if len(value) > 0 {
		parsedValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("failed to pares str to int: %s", err.Error()))
			return "", pkg.ErrParsingToInt
		}
		val = parsedValue
	}

	item, err := s.store.GetItem(key)
	if err != nil {
		return "", err
	}

	itemValue, err := strconv.ParseInt(item.Value, 10, 64)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("failed to pares str to int: %s", err.Error()))
		return "", pkg.ErrParsingToInt
	}

	switch command {
	case INC:
		itemValue += val
	case DEC:
		itemValue -= val
	}

	item.Value = strconv.Itoa(int(itemValue))

	return s.set(key, item.Value, item.TTL)
}

func (s *Server) inc(key string, value string) (string, error) {
	return s.updateIntValue(key, value, INC)
}

func (s *Server) dec(key string, value string) (string, error) {
	return s.updateIntValue(key, value, DEC)
}

func (s *Server) negate(key string) (string, error) {
	item, err := s.store.GetItem(key)
	if err != nil {
		return "", err
	}

	val, err := strconv.ParseBool(item.Value)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("failed to pares str to bool: %s", err.Error()))
		return "", pkg.ErrParsingToBool
	}

	item.Value = strconv.FormatBool(!val)

	return s.set(key, item.Value, item.TTL)
}

func (s *Server) updateList(key string, value string, command string) (string, error) {
	item, err := s.store.GetItem(key)
	if err != nil {
		return "", err
	}

	slice, err := pkg.ParseStringToSlice(item.Value)
	if err != nil {
		return "", pkg.ErrMalformedSlice
	}

	switch command {
	case LAPPEND:
		{
			slice = append(slice, value)
			item.Value = pkg.ParseSliceToString(slice)
		}
	case LREMOVE:
		{
			targetIndex := -1
			for i, val := range slice {
				if val == value {
					targetIndex = i
					break
				}
			}

			if targetIndex == -1 {
				return "", pkg.ErrElementNotinList
			}

			slice[targetIndex] = slice[len(slice)-1]
			item.Value = pkg.ParseSliceToString(slice[:len(slice)-1])
		}

	}

	return s.set(key, item.Value, item.TTL)
}

func (s *Server) lappend(key string, value string) (string, error) {
	return s.updateList(key, value, LAPPEND)
}

func (s *Server) lremove(key string, value string) (string, error) {
	return s.updateList(key, value, LREMOVE)
}
