package token

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	// types
	ID     = "ID"
	INT    = "INT"
	FLOAT  = "FLOAT"
	TYPE   = "TYPE"
	STRING = "STRING"

	// operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"
	POWER    = "^"
	EQ       = "=="
	NE       = "!="
	LT       = "<"
	GT       = ">"
	LE       = "<="
	GE       = ">="
	BANG     = "!"
	OR       = "||"
	AND      = "&&"
	DOT      = "."
	MOD      = "%"

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	// keywords
	IF     = "IF"
	ELSE   = "ELSE"
	FUNC   = "FUNC"
	LET    = "LET"
	RETURN = "RETURN"
	FALSE  = "FALSE"
	TRUE   = "TRUE"

	// others
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Row     int
	Column  int
}

func NewToken(t TokenType, lit byte) *Token {
	return &Token{
		Type:    t,
		Literal: string(lit),
	}
}

func NewTokenString(t TokenType, lit string) *Token {
	return &Token{
		Type:    t,
		Literal: lit,
	}
}

var keywords = map[string]TokenType{
	"if":     IF,
	"else":   ELSE,
	"fn":     FUNC,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

// Golang doesn't have sets, we use 0-sized
// structs as values instead
type void struct{}

var null void = void{}

var types = map[string]void{
	"Int":    null,
	"Float":  null,
	"Func":   null,
	"Void":   null,
	"Bool":   null,
	"String": null,
	"List":   null,
	"Map":    null,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, isKeyword := keywords[identifier]; isKeyword {
		return tokenType
	} else if _, isType := types[identifier]; isType {
		return TYPE
	} else {
		return ID
	}
}
