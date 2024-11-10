package server

type Lexer struct {
	input        string
	ch           byte
	readPosition int
	lastToken    *Token
	position     int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}

	l.readChar()

	return l
}

func (l *Lexer) NextToken() *Token {
	var token Token

	l.skipWhiteSpace()

	switch l.ch {
	case 0:
		token = NewToken(ILLEGAL, EOF)
	case '[':
		token = NewToken(LIST, l.readList())
	case '"':
		token = NewToken(STRING, l.readString())
	default:
		switch {
		case isChar(l.ch):
			literal := l.readId()
			if _, ok := commands[literal]; ok {
				token = NewToken(COMMAND, literal)
			} else if param, ok := parameters[literal]; ok {
				token = NewToken(param, literal)
			} else {
				token = NewToken(ID, literal)
			}
		case isDigit(l.ch):
			literal, mixedChar := l.readInt()
			var tokenType TokenType
			if mixedChar {
				tokenType = ID
			} else {
				tokenType = INT
			}
			if l.lastToken.Literal == TTL {
				tokenType = TTL
			}
			token = NewToken(tokenType, literal)
		default:
			token = NewToken(ILLEGAL, string(l.ch))
		}

	}

	l.lastToken = &token
	l.readChar()

	return &token
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readList() string {
	position := l.position

	for {
		l.readChar()
		if l.ch == ']' {
			break
		}
	}

	return l.input[position:l.readPosition]
}

func isChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}

func (l *Lexer) readInt() (string, bool) {
	position := l.position
	mixed := false
	for isDigit(l.ch) {
		l.readChar()
	}

	if isChar(l.ch) {
		mixed = true
		for isChar(l.ch) || isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position], mixed
}

func (l *Lexer) readId() string {
	position := l.position

	for isChar(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.readPosition
	var str string

	for {
		l.readChar()
		if l.ch == '"' {
			str = l.input[position:l.position]
			break
		}
	}

	return str
}

func (l *Lexer) skipWhiteSpace() {
	for isWhiteSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) stepBack() {
	if l.readPosition > 0 && l.position > 0 {
		l.ch = l.input[l.position]
		l.readPosition = l.position
		l.position -= 1
	}
}

func (l *Lexer) peekNextChar() byte {
	return l.input[l.readPosition]
}
