package lexer

import (
	"bufio"
	"io"
	"lang/token"
	"log"
	"strings"
)

type position struct {
	row    int
	column int
}

type Lexer struct {
	input *bufio.Reader

	char      byte
	curByte   int
	position  *position
	isNewline bool
}

func NewLexer(code string) *Lexer {
	reader := bufio.NewReader(strings.NewReader(code))
	l := &Lexer{
		input:    reader,
		position: &position{column: 1, row: 1},
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.curByte >= l.input.Size() {
		l.char = 0 // byte 0 is EOF
	}

	if l.isNewline {
		l.position.row++
		l.position.column = 1
		l.isNewline = false
	}

	l.char, _ = l.input.ReadByte()
	if l.char == '\n' {
		l.isNewline = true
	}

	l.curByte++
	l.position.column++
}

func (l *Lexer) peekChar() string {
	b, err := l.input.Peek(1)
	if err != nil && err != io.EOF {
		log.Fatalf("%q error peeking while Lexing %v", l.position, err)
	} else if err == io.EOF {
		return string([]byte{0})
	}
	return string(b)
}

func (l *Lexer) readIdentifier() string {
	var out []byte
	out = append(out, l.char)

	for isAlphanumeric([]byte(l.peekChar())[0]) {
		l.readChar()
		out = append(out, l.char)
	}

	return string(out)
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	var out []byte
	var t token.TokenType = token.INT
	var expectFloat bool = true

	out = append(out, l.char)

	for isDigit([]byte(l.peekChar())[0], expectFloat) {
		if l.char == '.' {
			t = token.FLOAT
			expectFloat = false
		}

		l.readChar()
		out = append(out, l.char)
	}

	return string(out), t
}

func (l *Lexer) readString() string {
	var out []byte
	l.readChar()
	for {
		if l.char == '"' || l.char == 0 {
			break
		}
		out = append(out, l.char)
		l.readChar()
	}

	return string(out)
}

func (l *Lexer) NextToken() *token.Token {
	var t *token.Token
	l.skipWhitespaces()

	switch l.char {
	case '"':
		str := l.readString()
		t = token.NewTokenString(token.STRING, str)
	case '=':
		if l.isPeek('=') {
			l.readChar()
			t = token.NewTokenString(token.EQ, "==")
		} else {
			t = token.NewToken(token.ASSIGN, l.char)
		}
	case '+':
		t = token.NewToken(token.PLUS, l.char)
	case '-':
		t = token.NewToken(token.MINUS, l.char)
	case '*':
		t = token.NewToken(token.ASTERISK, l.char)
	case '/':
		if l.isPeek('/') {
			l.skipComment()
			return l.NextToken()
		} else {
			t = token.NewToken(token.SLASH, l.char)
		}
	case '^':
		t = token.NewToken(token.POWER, l.char)
	case '>':
		if l.isPeek('=') {
			l.readChar()
			t = token.NewTokenString(token.GE, ">=")
		} else {
			t = token.NewToken(token.GT, l.char)
		}
	case '<':
		if l.isPeek('=') {
			l.readChar()
			t = token.NewTokenString(token.LE, "<=")
		} else {
			t = token.NewToken(token.LT, l.char)
		}
	case '!':
		if l.isPeek('=') {
			l.readChar()
			t = token.NewTokenString(token.NE, "!=")
		} else {
			t = token.NewToken(token.BANG, l.char)
		}
	case '|':
		l.readChar()
		t = token.NewTokenString(token.OR, "||")
	case '&':
		l.readChar()
		t = token.NewTokenString(token.AND, "&&")
	case '%':
		t = token.NewToken(token.MOD, l.char)
	case '(':
		t = token.NewToken(token.LPAREN, l.char)
	case ')':
		t = token.NewToken(token.RPAREN, l.char)
	case '{':
		t = token.NewToken(token.LBRACE, l.char)
	case '}':
		t = token.NewToken(token.RBRACE, l.char)
	case '[':
		t = token.NewToken(token.LBRACKET, l.char)
	case ']':
		t = token.NewToken(token.RBRACKET, l.char)
	case ';':
		t = token.NewToken(token.SEMICOLON, l.char)
	case ',':
		t = token.NewToken(token.COMMA, l.char)
	case ':':
		t = token.NewToken(token.COLON, l.char)
	case '.':
		if isDigit(l.peekChar()[0], false) {
			numberLiteral, tokenType := l.readNumber()
			t = token.NewTokenString(tokenType, numberLiteral)
		} else {
			t = token.NewToken(token.DOT, '.')
		}
	case 0:
		t = token.NewToken(token.EOF, l.char)
	default:
		if isLetter(l.char) {
			identifier := l.readIdentifier()
			t = token.NewTokenString(token.LookupIdentifier(identifier), identifier)
		} else if isDigit(l.char, true) {
			numberLiteral, tokenType := l.readNumber()
			t = token.NewTokenString(tokenType, numberLiteral)
		} else {
			t = token.NewToken(token.ILLEGAL, l.char)
		}
	}

	l.readChar()
	l.setPosition(t)
	return t
}

func (l *Lexer) isPeek(char byte) bool {
	return string(char) == l.peekChar()
}

func (l *Lexer) setPosition(t *token.Token) {
	t.Column = l.position.column - len(t.Literal) - 1
	t.Row = l.position.row
}

func (l *Lexer) skipWhitespaces() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func isLetter(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_'
}

func isDigit(char byte, expectFloat bool) bool {
	if expectFloat {
		return char >= '0' && char <= '9' || char == '.'
	} else {
		return char >= '0' && char <= '9'
	}
}

func isAlphanumeric(char byte) bool {
	return isLetter(char) || isDigit(char, false)
}

func (l *Lexer) skipComment() {
	for l.char != '\n' {
		l.readChar()
	}
}
