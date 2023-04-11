package parser

import (
	"fmt"
	"lang/ast"

	"github.com/m1gwings/treedrawer/tree"
)

func DrawTree(p *ast.Program) {
	t := tree.NewTree(tree.NodeString("Program"))
	for _, stmt := range p.Statements {
		drawStatement(stmt, t)
	}
	fmt.Println(t)
}

func drawStatement(stmt ast.Statement, parent *tree.Tree) {
	switch stmt := stmt.(type) {
	case *ast.LetStatement:
		child := parent.AddChild(tree.NodeString(stmt.TokenLiteral()))
		drawLet(stmt, child)
	case *ast.ReturnStatement:
		child := parent.AddChild(tree.NodeString("return"))
		drawReturn(stmt, child)
	case *ast.BlockStatement:
		drawBlockStatement(stmt, parent)
	case *ast.ExpressionStatement:
		drawExpression(stmt.Expression, parent)
	}
}

func drawLet(stmt *ast.LetStatement, parent *tree.Tree) {
	parent.AddChild(tree.NodeString(stmt.Name.String()))
	drawExpression(stmt.Value, parent)
}

func drawExpression(exp ast.Expression, parent *tree.Tree) {
	switch exp := exp.(type) {
	case *ast.Identifier:
		drawIdentifier(exp, parent)
	case *ast.IntegerLiteral:
		drawIntegerLiteral(exp, parent)
	case *ast.FloatLiteral:
		drawFloatLiteral(exp, parent)
	case *ast.StringLiteral:
		drawStringLiteral(exp, parent)
	case *ast.InfixExpression:
		drawInfixExpression(exp, parent)
	case *ast.PrefixExpression:
		drawPrefixExpression(exp, parent)
	case *ast.CallExpression:
		drawCallExpression(exp, parent)
	case *ast.IfExpression:
		drawIfExpression(exp, parent)
	case *ast.Function:
		drawFunction(exp, parent)
	case *ast.FunctionLiteral:
		drawFunctionLiteral(exp, parent)
	case *ast.ListLiteral:
		child := parent.AddChild(tree.NodeString("list"))
		drawListLiteral(exp, child)
	case *ast.AccessExpression:
		drawAccessExpression(exp, parent)
	}
}

func drawIdentifier(exp *ast.Identifier, parent *tree.Tree) {
	parent.AddChild(tree.NodeString(exp.Value))
}

func drawIntegerLiteral(exp *ast.IntegerLiteral, parent *tree.Tree) {
	parent.AddChild(tree.NodeInt64(exp.Value))
}

func drawStringLiteral(exp *ast.StringLiteral, parent *tree.Tree) {
	parent.AddChild(tree.NodeString("\"" + exp.Value + "\""))
}

func drawFloatLiteral(exp *ast.FloatLiteral, parent *tree.Tree) {
	parent.AddChild(tree.NodeFloat64(exp.Value))
}

func drawInfixExpression(exp *ast.InfixExpression, parent *tree.Tree) {
	child := parent.AddChild(tree.NodeString(exp.Operator))
	drawExpression(exp.Left, child)
	drawExpression(exp.Right, child)
}

func drawPrefixExpression(exp *ast.PrefixExpression, parent *tree.Tree) {
	child := parent.AddChild(tree.NodeString(exp.Operator))
	drawExpression(exp.Right, child)
}

func drawCallExpression(exp *ast.CallExpression, parent *tree.Tree) {
	child := parent.AddChild(tree.NodeString("call"))
	drawExpression(exp.Function, child)
	for _, args := range exp.Arguments {
		drawExpression(args, child)
	}
}

func drawBlockStatement(block *ast.BlockStatement, parent *tree.Tree) {
	for _, s := range block.Statements {
		drawStatement(s, parent)
	}
}

func drawIfExpression(exp *ast.IfExpression, parent *tree.Tree) {
	child := parent.AddChild(tree.NodeString(exp.TokenLiteral()))
	drawExpression(exp.Condition, child)
	thenNode := child.AddChild(tree.NodeString("then"))
	drawStatement(exp.Consequence, thenNode)
	if exp.Alternative != nil {
		elseNode := child.AddChild(tree.NodeString("else"))
		drawStatement(exp.Alternative, elseNode)
	}
}

func drawFunction(fn *ast.Function, parent *tree.Tree) {
	child := parent.AddChild("define " + tree.NodeString(fn.Name.Value))
	nextChild := addParameters(fn.FunctionLiteral, child)
	drawBlockStatement(fn.Body, nextChild)
}

func addParameters(fn *ast.FunctionLiteral, parent *tree.Tree) *tree.Tree {
	if len(fn.Parameters) == 0 {
		return parent
	}
	params := ""
	pr := fn.Parameters
	for _, param := range pr {
		params += param.ParamString()
	}
	return parent.AddChild(tree.NodeString(params))
}

func drawFunctionLiteral(fn *ast.FunctionLiteral, parent *tree.Tree) {
	child := parent.AddChild(tree.NodeString("define"))
	childNode := addParameters(fn, child)
	drawBlockStatement(fn.Body, childNode)
}

func drawReturn(stmt *ast.ReturnStatement, parent *tree.Tree) {
	drawExpression(stmt.ReturnValue, parent)
}

func drawListLiteral(list *ast.ListLiteral, parent *tree.Tree) {
	for _, value := range list.Elements {
		drawExpression(value, parent)
	}
}

func drawAccessExpression(exp *ast.AccessExpression, parent *tree.Tree) {
	drawExpression(exp.Struct, parent)
	parent.AddChild(tree.NodeString(exp.Attribute))
}
