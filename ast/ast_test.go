package ast

import (
	"lang/token"
	"testing"
)

func TestString(t *testing.T) {
	expected := "let Var: int = anotherVar;"
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.ID, Literal: "Var"},
					Value: "Var",
					Type:  token.Token{Type: token.TYPE, Literal: "int"},
				},
				Value: &Identifier{
					Token: token.Token{Type: token.ID, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != expected {
		t.Errorf("program.String() expected %s, got=%s", expected, program.String())
	}
}
