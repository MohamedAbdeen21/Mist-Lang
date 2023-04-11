package ast

import (
	"bytes"
	"lang/token"
	"strings"
)

type Node interface {
	// mainly for testing and debugging
	String() string
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

// every expression can return a value
type Expression interface {
	Node
	ReturnType() string
	expressionNode()
}

// root of every ast
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token // token.ID
	Type  token.Token // token.TYPE
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) ReturnType() string   { return i.Type.Literal }
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) ParamString() string  { return i.Value + ": " + i.Type.Literal }

type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(": " + ls.Name.Type.Literal)
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token // token.INT
	Type  token.Token // token.TYPE
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) ReturnType() string   { return il.Type.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token // token.FLOAT
	Type  token.Token // token.TYPE
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) ReturnType() string   { return fl.Type.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Type     token.Token // token.TYPE
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) ReturnType() string   { return pe.Type.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator + pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Type     token.Token // token.TYPE
	Operator string
	Right    Expression
	Left     Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) ReturnType() string   { return ie.Type.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String() + " " + ie.Operator + " " + ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Type  token.Token // token.TYPE
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) ReturnType() string   { return b.Type.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token // token.IF
	Type        token.Token // token.TYPE
	Condition   Expression
	Consequence *BlockStatement
	Others      Expression // else if
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) ReturnType() string   { return ie.Type.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" " + ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else " + ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // token.FUNC
	Type       token.Token // token.TYPE
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) ReturnType() string   { return fl.Type.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.ParamString())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(" + strings.Join(params, ", ") + ") ")
	out.WriteString(fl.Type.Literal + " ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // token.LPAREN
	Type      token.Token // token.TYPE
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) ReturnType() string   { return ce.Type.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(" + strings.Join(args, ", ") + ")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token // token.STRING
	Type  token.Token // token.TYPE
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) ReturnType() string   { return sl.Type.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type Function struct {
	*FunctionLiteral
	Name *Identifier
}

func (f *Function) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.ParamString())
	}

	out.WriteString(f.TokenLiteral() + " " + f.Name.Value)
	out.WriteString("(" + strings.Join(params, ", ") + ") ")
	out.WriteString(f.Type.Literal + " ")
	out.WriteString(f.Body.String())

	return out.String()
}

type ListLiteral struct {
	Token    token.Token // token.LBRACKET
	Type     token.Token // token.TYPE
	Elements []Expression
}

func (ll *ListLiteral) expressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }
func (ll *ListLiteral) ReturnType() string   { return ll.Type.Literal }
func (ll *ListLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range ll.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[" + strings.Join(elements, ", ") + "]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token // token.LBRACKET
	Type  token.Token // token.TYPE
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) ReturnType() string   { return ie.Type.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[" + ie.Index.String() + "]")
	out.WriteString(")")

	return out.String()
}

type AccessExpression struct {
	Token     token.Token // token.DOT
	Type      token.Token // token.TYPE
	Struct    Expression
	Attribute string
}

func (ae *AccessExpression) expressionNode()      {}
func (ae *AccessExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AccessExpression) ReturnType() string   { return ae.Type.Literal }

func (ae *AccessExpression) String() string { return ae.Struct.String() + "." + ae.Attribute }

type MapLiteral struct {
	Token     token.Token // token.LBRACE
	Type      token.Token // token.TYPE
	KeyType   token.Token // token.TYPE
	ValueType token.Token // token.TYPE
	Pairs     map[Expression]Expression
}

func (ml *MapLiteral) expressionNode()      {}
func (ml *MapLiteral) TokenLiteral() string { return ml.Token.Literal }

func (ml *MapLiteral) ReturnType() string { return ml.Type.Literal }
func (ml *MapLiteral) String() string {
	pairs := []string{}
	for key, value := range ml.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	return "{" + strings.Join(pairs, ", ") + "}"
}
