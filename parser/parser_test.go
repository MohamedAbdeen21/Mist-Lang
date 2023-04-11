package parser

import (
	"lang/ast"
	"lang/lexer"
	"testing"
)

func checkParserErrors(index int, t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("case %d: parser has %d errors", index, len(errors))
	for _, msg := range errors {
		t.Errorf("--- parse error: %q", msg)
	}

	t.FailNow()
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		code        string
		id          string
		expression  string
		typeLiteral string
	}{
		{code: "let x: Int = 5;", id: "x", expression: "5", typeLiteral: "Int"},
		{code: "let y: Float = .75;", id: "y", expression: ".75", typeLiteral: "Float"},
		{code: "let z: Float = y;", id: "z", expression: "y", typeLiteral: "Float"},
		{code: "let a: Func = b;", id: "a", expression: "b", typeLiteral: "Func"},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if program == nil {
			t.Fatalf("Parse() returned nil")
		}

		stmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Fatalf("case %d: expected program.Statements[0] to be ast.LetStatement, got=%T",
				i, program.Statements[0])
		}

		if stmt.Name.ReturnType() != test.typeLiteral {
			t.Errorf("case %d: expected type %s, got=%s",
				i, stmt.Name.ReturnType(), test.typeLiteral)
		}

		if stmt.Value.ReturnType() != "" {
			if stmt.Name.ReturnType() != stmt.Value.ReturnType() {
				t.Errorf("case %d: expected type of value %s, to match type of id, got=%s",
					i, stmt.Name.ReturnType(), stmt.Value.ReturnType())
			}
		}

		if stmt.Name.Value != test.id {
			t.Errorf("case %d: expected literal %s, got=%s",
				i, stmt.Token.Literal, test.id)
		}

		// if stmt.Value.TokenLiteral() != test.expression {
		// 	t.Errorf("case %d: expected value %s, got=%s",
		// 		i, test.expression, stmt.Value.TokenLiteral())
		// }

	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		code         string
		tokenLiteral string
		expression   string
	}{
		{code: "return 5;", tokenLiteral: "return"},
		{code: "return .5;", tokenLiteral: "return"},
		{code: "return x;", tokenLiteral: "return"},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("case %d: expected program.Statements[0] to be ast.ReturnStatement, got=%T",
				i, program.Statements[0])
		}

		if stmt.Token.Literal != test.tokenLiteral {
			t.Errorf("case %d: expected literal %s, got=%s",
				i, stmt.Token.Literal, test.tokenLiteral)
		}

		// if stmt.Value.TokenLiteral() != test.expression {
		// 	t.Errorf("case %d: expected value %s, got=%s",
		// 		i, test.expression, stmt.Value.TokenLiteral())
		// }

	}
}

func TestIdentifierExpression(t *testing.T) {
	tests := []struct {
		code         string
		tokenLiteral string
		expression   string
	}{
		{code: "foobar;", tokenLiteral: "foobar", expression: "foobar"},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"case %d: expected program.Statements[0] to be ast.ExpressionStatement, got=%T",
				i,
				program.Statements[0],
			)
		}

		ident, ok := stmt.Expression.(*ast.Identifier)
		if !ok {
			t.Fatalf("expression not *ast.Identifier, got=%T", stmt.Expression)
		}

		if ident.Value != test.expression {
			t.Errorf("ident.Value not %s, got=%s", test.expression, ident.Value)
		}

		if ident.TokenLiteral() != test.tokenLiteral {
			t.Errorf("ident.TokenLiteral() not %s, got=%s", test.tokenLiteral, ident.TokenLiteral())
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	tests := []struct {
		code         string
		tokenLiteral string
		expression   int64
	}{
		{code: "5", tokenLiteral: "5", expression: 5},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"case %d: expected program.Statements[0] to be ast.ExpressionStatement, got=%T",
				i,
				program.Statements[0],
			)
		}

		ident, ok := stmt.Expression.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("expression not *ast.IntergerLiteral, got=%T", stmt.Expression)
		}

		if ident.Value != test.expression {
			t.Errorf("ident.Value not %d, got=%d", test.expression, ident.Value)
		}

		if ident.TokenLiteral() != test.tokenLiteral {
			t.Errorf("ident.TokenLiteral() not %s, got=%s", test.tokenLiteral, ident.TokenLiteral())
		}

		// if stmt.Value.TokenLiteral() != test.expression {
		// 	t.Errorf("case %d: expected value %s, got=%s",
		// 		i, test.expression, stmt.Value.TokenLiteral())
		// }

	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		code     string
		operator string
		value    string
	}{
		{code: "-5", operator: "-", value: "5"},
		{code: "!15", operator: "!", value: "15"},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"case %d: expected program.Statements[0] to be ast.ExpressionStatement, got=%T",
				i,
				program.Statements[0],
			)
		}

		prefix, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expression not *ast.PrefixExpression, got=%T", stmt.Expression)
		}

		if prefix.Operator != test.operator {
			t.Errorf("prefix.Operator not %s, got=%s", test.operator, prefix.Operator)
		}

		if prefix.Right.String() != test.value {
			t.Errorf("preifx.Right not %s, got=%s", test.value, prefix.Right.String())
		}

		// if stmt.Value.TokenLiteral() != test.expression {
		// 	t.Errorf("case %d: expected value %s, got=%s",
		// 		i, test.expression, stmt.Value.TokenLiteral())
		// }

	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		code       string
		operator   string
		leftValue  string
		rightValue string
	}{
		{code: "5+6;", operator: "+", leftValue: "5", rightValue: "6"},
		{code: "10*12;", operator: "*", leftValue: "10", rightValue: "12"},
		{code: "2/1;", operator: "/", leftValue: "2", rightValue: "1"},
		{code: "5-3;", operator: "-", leftValue: "5", rightValue: "3"},
		{code: "5-3;", operator: "-", leftValue: "5", rightValue: "3"},
		{code: "5<3;", operator: "<", leftValue: "5", rightValue: "3"},
		{code: "5>11;", operator: ">", leftValue: "5", rightValue: "11"},
		{code: "5==11;", operator: "==", leftValue: "5", rightValue: "11"},
		{code: "5!=11;", operator: "!=", leftValue: "5", rightValue: "11"},
		{code: "5||11;", operator: "||", leftValue: "5", rightValue: "11"},
		{code: "5&&11;", operator: "&&", leftValue: "5", rightValue: "11"},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"case %d: expected program.Statements[0] to be ast.ExpressionStatement, got=%T",
				i,
				program.Statements[0],
			)
		}

		infix, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expression not *ast.InfixExpression, got=%T", stmt.Expression)
		}

		if infix.Operator != test.operator {
			t.Errorf("infix.Operator not %s, got=%s", test.operator, infix.Operator)
		}

		if infix.Right.String() != test.rightValue {
			t.Errorf("infix.Right not %s, got=%s", test.rightValue, infix.Right)
		}

		if infix.Left.String() != test.leftValue {
			t.Errorf("infix.Left not %s, got=%s", test.leftValue, infix.Right)
		}

		// if stmt.Value.TokenLiteral() != test.expression {
		// 	t.Errorf("case %d: expected value %s, got=%s",
		// 		i, test.expression, stmt.Value.TokenLiteral())
		// }

	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{code: "-a+b", expected: "((-a) + b)"},
		{code: "!-a", expected: "(!(-a))"},
		{code: "a + b + c", expected: "((a + b) + c)"},
		{code: "a + b * c / d", expected: "(a + ((b * c) / d))"},
		{code: "a + b * (c / d)", expected: "(a + (b * (c / d)))"},
		{code: "5 > 4 == 3 < 4", expected: "((5 > 4) == (3 < 4))"},
		{code: "5 || 3 == 4", expected: "(5 || (3 == 4))"},
		{code: "5 || (3 && 4)", expected: "(5 || (3 && 4))"},
		{code: "5 || add(3,4) * 5", expected: "(5 || (add(3, 4) * 5))"},
		{code: "5 <= 3 >= 4", expected: "((5 <= 3) >= 4)"},
		{code: "(3 + 5) * 4 != 3 * 1 + 24", expected: "(((3 + 5) * 4) != ((3 * 1) + 24))"},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		actual := program.String()
		if actual != test.expected {
			t.Errorf("case %d: expected=%s, got=%s", i, test.expected, actual)
		}
	}
}

func TestBoolean(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{code: "!true", expected: "(!true)"},
		{code: "false + true", expected: "(false + true)"},
		{code: "true", expected: "true"},
		{code: "false", expected: "false"},
		{code: "3 > 5 == true", expected: "((3 > 5) == true)"},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		actual := program.String()
		if actual != test.expected {
			t.Errorf("case %d: expected=%s, got=%s", i, test.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		code                string
		expectedCondition   string
		expectedConsequence string
		expectedAlternative string
	}{
		{
			code:                "if (x < y) {x}",
			expectedCondition:   "(x < y)",
			expectedConsequence: "x",
		},
		{
			code:                "if (x < y) {x} else {y}",
			expectedCondition:   "(x < y)",
			expectedConsequence: "x",
			expectedAlternative: "y",
		},
		{
			code:                "if (x == y) {x + 1} else {y - 1}",
			expectedCondition:   "(x == y)",
			expectedConsequence: "(x + 1)",
			expectedAlternative: "(y - 1)",
		},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0])
		}

		cond, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression, got=%T", stmt.Expression)
		}

		if cond.Condition.String() != test.expectedCondition {
			t.Errorf("expected condition to be %s, got=%s",
				test.expectedCondition, cond.Condition.String())
		}

		if cond.Consequence.String() != test.expectedConsequence {
			t.Errorf("expected consequence to be %s, got=%s",
				test.expectedCondition, cond.Condition.String())
		}

		if test.expectedAlternative != "" {
			if cond.Alternative == nil {
				t.Errorf("expected Alternative to be %s, got nil", test.expectedAlternative)
			}
			if cond.Alternative.String() != test.expectedAlternative {
				t.Errorf("expected Alternative to be %s, got=%s",
					test.expectedCondition, cond.Condition.String())
			}
		}
	}
}

func TestFunctionLiteral(t *testing.T) {
	tests := []struct {
		code               string
		expectedBody       string
		expectedParameters []string
		expectedReturnType string
	}{
		{
			code:               "fn(x:Int,y:Int) Int {x+y;}",
			expectedBody:       "(x + y)",
			expectedParameters: []string{"x: Int", "y: Int"},
			expectedReturnType: "Int",
		},
		{
			code:               "fn(x:Float,y:Float) Float {x*y;}",
			expectedBody:       "(x * y)",
			expectedParameters: []string{"x: Float", "y: Float"},
			expectedReturnType: "Float",
		},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0])
		}

		fn, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.FunctionLiteral, got=%T",
				program.Statements[0])
		}

		if fn.Body.String() != test.expectedBody {
			t.Errorf("expected body to be %s, got=%s", test.expectedBody, fn.Body.String())
		}

		for i, param := range fn.Parameters {
			if test.expectedParameters[i] != param.ParamString() {
				t.Errorf("expected parameter %d to be %s, got=%s", i,
					test.expectedParameters[i], param.ParamString())
			}
		}

		if fn.ReturnType() != test.expectedReturnType {
			t.Errorf("expected return type to be %s, got=%s",
				test.expectedReturnType, fn.ReturnType())
		}
	}
}

func TestCallExpression(t *testing.T) {
	tests := []struct {
		code              string
		expectedFunction  string
		expectedArguments []string
	}{
		{
			code:              "add(3,4);",
			expectedFunction:  "add",
			expectedArguments: []string{"3", "4"},
		},
		{
			code:              "sum(3,5,1)",
			expectedFunction:  "sum",
			expectedArguments: []string{"3", "5", "1"},
		},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0])
		}

		call, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression, got=%T",
				program.Statements[0])
		}

		if call.Function.String() != test.expectedFunction {
			t.Errorf("expected body to be %s, got=%s",
				test.expectedFunction, call.Function.String())
		}

		for i, arg := range call.Arguments {
			if test.expectedArguments[i] != arg.String() {
				t.Errorf("expected parameter %d to be %s, got=%s", i,
					test.expectedArguments[i], arg.String())
			}
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	tests := []struct {
		code         string
		tokenLiteral string
		expression   int64
	}{
		{code: "[1,2*2,3+3]", tokenLiteral: "5", expression: 5},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"case %d: expected program.Statements[0] to be ast.ExpressionStatement, got=%T",
				i,
				program.Statements[0],
			)
		}

		array, ok := stmt.Expression.(*ast.ListLiteral)
		if !ok {
			t.Fatalf("expression not *ast.ArrayLiteral, got=%T", stmt.Expression)
		}

		if len(array.Elements) != 3 {
			t.Errorf("len(array.Elements) not 3, got=%d", len(array.Elements))
		}
	}
}

func TestMapLiter(t *testing.T) {
	tests := []struct {
		code     string
		expected map[string]int64
	}{
		{code: `{"one":1,"two":2}`, expected: map[string]int64{"one": 1, "two": 2, "three": 3}},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.code)
		p := NewParser(l)

		program := p.Parse()
		checkParserErrors(i, t, p)

		if len(program.Statements) == 0 {
			t.Fatalf("expected program.Statements to have %d Statement(s), got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"case %d: expected program.Statements[0] to be ast.ExpressionStatement, got=%T",
				i,
				program.Statements[0],
			)
		}

		mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
		if !ok {
			t.Fatalf("expression not *ast.MapLiteral, got=%T", stmt.Expression)
		}

		for key, value := range mapLiteral.Pairs {
			literal, ok := key.(*ast.StringLiteral)
			if !ok {
				t.Errorf("key is not ast.StringLiteral, got=%T", key)
			}

			expectedValue, ok := test.expected[literal.String()]
			if !ok {
				t.Fatalf("Key is invalid, got=%s", literal.String())
			}

			integer, ok := value.(*ast.IntegerLiteral)
			if !ok {
				t.Errorf("value is not ast.IntegerLiteral, got=%T", value)
			}

			if expectedValue != integer.Value {
				t.Fatalf("Value should be %d, got=%d", expectedValue, integer.Value)
			}
		}
	}
}
