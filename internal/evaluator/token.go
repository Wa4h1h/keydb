package evaluator

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// operations
	GET     = "GET"
	SET     = "SET"
	INC     = "INC"
	DEC     = "DEC"
	NEGATE  = "NEGATE"
	LAPPEND = "LAPPEND"
	LREMOVE = "LREMOVE"
	REMOVE  = "REMOVE"

	// parameters
	TTL = "ttl"

	// type values
	STRING = "STRING"
	LIST   = "LIST"
	INT    = "INT"
	EOF    = "EOF"

	// identifier
	COMMAND = "COMMAND"
	ID      = "ID"
	ILLEGAL = "ILLEGAL"
)

var commands = map[string]TokenType{
	"GET":     GET,
	"SET":     SET,
	"REMOVE":  REMOVE,
	"INC":     INC,
	"DEC":     DEC,
	"NEGATE":  NEGATE,
	"LAPPEND": LAPPEND,
	"LREMOVE": LREMOVE,
}

var parameters = map[string]TokenType{
	"TTL": TTL,
}

func NewToken(t TokenType, l string) Token {
	return Token{Type: t, Literal: l}
}
