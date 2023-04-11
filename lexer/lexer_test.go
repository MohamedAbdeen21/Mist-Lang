package lexer

import (
	"lang/token"
	"testing"
)

func TestLexer(t *testing.T) {
	input := []struct {
		code     string
		expected []token.Token
	}{
		{
			"let _name: Int = 5;",
			[]token.Token{
				{Type: token.LET, Literal: "let", Row: 1, Column: 1},
				{Type: token.ID, Literal: "_name", Row: 1, Column: 5},
				{Type: token.COLON, Literal: ":", Row: 1, Column: 10},
				{Type: token.TYPE, Literal: "Int", Row: 1, Column: 12},
				{Type: token.ASSIGN, Literal: "=", Row: 1, Column: 16},
				{Type: token.INT, Literal: "5", Row: 1, Column: 18},
				{Type: token.SEMICOLON, Literal: ";", Row: 1, Column: 19},
				{Type: token.EOF, Literal: "\x00", Row: 1, Column: 20},
			},
		},

		{
			"let _name = 5;\nlet first_name_1: Float = .72;\t",
			[]token.Token{
				{Type: token.LET, Literal: "let", Row: 1, Column: 1},
				{Type: token.ID, Literal: "_name", Row: 1, Column: 5},
				{Type: token.ASSIGN, Literal: "=", Row: 1, Column: 11},
				{Type: token.INT, Literal: "5", Row: 1, Column: 13},
				{Type: token.SEMICOLON, Literal: ";", Row: 1, Column: 14},
				{Type: token.LET, Literal: "let", Row: 2, Column: 1},
				{Type: token.ID, Literal: "first_name_1", Row: 2, Column: 5},
				{Type: token.COLON, Literal: ":", Row: 2, Column: 17},
				{Type: token.TYPE, Literal: "Float", Row: 2, Column: 19},
				{Type: token.ASSIGN, Literal: "=", Row: 2, Column: 25},
				{Type: token.FLOAT, Literal: ".72", Row: 2, Column: 27},
				{Type: token.SEMICOLON, Literal: ";", Row: 2, Column: 30},
				{Type: token.EOF, Literal: "\x00", Row: 2, Column: 32},
			},
		},

		{
			"5 > 1.2; \"some string\"[1,2]",
			[]token.Token{
				{Type: token.INT, Literal: "5", Row: 1, Column: 1},
				{Type: token.GT, Literal: ">", Row: 1, Column: 3},
				{Type: token.FLOAT, Literal: "1.2", Row: 1, Column: 5},
				{Type: token.SEMICOLON, Literal: ";", Row: 1, Column: 8},
				{Type: token.STRING, Literal: "some string", Row: 1, Column: 12},
				{Type: token.LBRACKET, Literal: "[", Row: 1, Column: 23},
				{Type: token.INT, Literal: "1", Row: 1, Column: 24},
				{Type: token.COMMA, Literal: ",", Row: 1, Column: 25},
				{Type: token.INT, Literal: "2", Row: 1, Column: 26},
				{Type: token.RBRACKET, Literal: "]", Row: 1, Column: 27},
				{Type: token.EOF, Literal: "\x00", Row: 1, Column: 28},
			},
		},

		{
			"=(){}+-/*let fn==<=>=<!>!=false;,true||return&&else:if Func",
			[]token.Token{
				{Type: token.ASSIGN, Literal: "=", Row: 1, Column: 1},
				{Type: token.LPAREN, Literal: "(", Row: 1, Column: 2},
				{Type: token.RPAREN, Literal: ")", Row: 1, Column: 3},
				{Type: token.LBRACE, Literal: "{", Row: 1, Column: 4},
				{Type: token.RBRACE, Literal: "}", Row: 1, Column: 5},
				{Type: token.PLUS, Literal: "+", Row: 1, Column: 6},
				{Type: token.MINUS, Literal: "-", Row: 1, Column: 7},
				{Type: token.SLASH, Literal: "/", Row: 1, Column: 8},
				{Type: token.ASTERISK, Literal: "*", Row: 1, Column: 9},
				{Type: token.LET, Literal: "let", Row: 1, Column: 10},
				{Type: token.FUNC, Literal: "fn", Row: 1, Column: 14},
				{Type: token.EQ, Literal: "==", Row: 1, Column: 16},
				{Type: token.LE, Literal: "<=", Row: 1, Column: 18},
				{Type: token.GE, Literal: ">=", Row: 1, Column: 20},
				{Type: token.LT, Literal: "<", Row: 1, Column: 22},
				{Type: token.BANG, Literal: "!", Row: 1, Column: 23},
				{Type: token.GT, Literal: ">", Row: 1, Column: 24},
				{Type: token.NE, Literal: "!=", Row: 1, Column: 25},
				{Type: token.FALSE, Literal: "false", Row: 1, Column: 27},
				{Type: token.SEMICOLON, Literal: ";", Row: 1, Column: 32},
				{Type: token.COMMA, Literal: ",", Row: 1, Column: 33},
				{Type: token.TRUE, Literal: "true", Row: 1, Column: 34},
				{Type: token.OR, Literal: "||", Row: 1, Column: 38},
				{Type: token.RETURN, Literal: "return", Row: 1, Column: 40},
				{Type: token.AND, Literal: "&&", Row: 1, Column: 46},
				{Type: token.ELSE, Literal: "else", Row: 1, Column: 48},
				{Type: token.COLON, Literal: ":", Row: 1, Column: 52},
				{Type: token.IF, Literal: "if", Row: 1, Column: 53},
				{Type: token.TYPE, Literal: "Func", Row: 1, Column: 56},
				{Type: token.EOF, Literal: "\x00", Row: 1, Column: 60},
			},
		},
	}

	for i, test := range input {
		l := NewLexer(test.code)
		for j, expected := range test.expected {
			actual := l.NextToken()
			if expected != *actual {
				t.Errorf("case %d:token %d expected %#v, got=%#v", i, j, expected, actual)
			}
		}
	}
}
