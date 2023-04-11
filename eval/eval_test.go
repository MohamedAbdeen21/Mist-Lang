package eval

import (
	"lang/lexer"
	"lang/object"
	"lang/parser"
	"testing"
)

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Parse()
	return Eval(program, object.NewScope())
}

func testInterface(t *testing.T, i int, expected interface{}, obj object.Object) {
	if integer, ok := expected.(int); ok {
		testIntegerObject(t, i, obj, int64(integer))
	} else if float, ok := expected.(float64); ok {
		testFloatObject(t, i, obj, float)
	} else if str, ok := expected.(string); ok {
		testStringObject(t, i, obj, str)
	} else if boolean, ok := expected.(bool); ok {
		testBooleanObject(t, i, obj, boolean)
	} else {
		testNullObject(t, i, obj)
	}
}

func testIntegerObject(t *testing.T, i int, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("case %d: object is not an integer, got=%T (%+v)", i, obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("case %d: expected integer object to have %d, got=%d", i, expected, result.Value)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, i int, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("case %d: object is not a float, got=%T (%+v)", i, obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("case %d: expected float object to have %f, got=%f", i, expected, result.Value)
		return false
	}

	return true
}

func testStringObject(t *testing.T, i int, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("case %d: object is not a string, got=%T (%+v)", i, obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("case %d: expected string object to have %s, got=%s", i, expected, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, i int, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("case %d: object is not a boolean, got=%T (%+v)", i, obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("case %d: expected boolean object to have %t, got=%t", i, expected, result.Value)
		return false
	}

	return true
}

func testNullObject(t *testing.T, i int, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("case %d: object is not NULL, got=%T (%+v)", i, obj, obj)
		return false
	}
	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		code     string
		expected int64
	}{
		{code: "5", expected: 5},
		{code: "10", expected: 10},
		{code: "-10", expected: -10},
		{code: "10 + 30", expected: 40},
		{code: "10 + 30 * 2", expected: 70},
		{code: "30 * 2 / 2 + 10", expected: 40},
		{code: "(1+1)^10 - 1 * 4", expected: 1020},
		{code: "20 + -10", expected: 10},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		testIntegerObject(t, i, evaluated, test.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		code     string
		expected bool
	}{
		{code: "true", expected: true},
		{code: "!true", expected: false},
		{code: "!!true", expected: true},
		{code: "false", expected: false},
		{code: "false != true", expected: true},
		{code: "false != false", expected: false},
		{code: "false || true", expected: true},
		{code: "false || false", expected: false},
		{code: "!false || false", expected: true},
		{code: "!false && false", expected: false},
		{code: "!false && true", expected: true},
		{code: "1 == 1", expected: true},
		{code: "1^(-1*-1) == 4*1/4", expected: true},
		{code: "8 >= 2 ^ 3 + 1", expected: false},
		{code: "\"Hello\" != \"world!\"", expected: true},
		{code: "\"Hello\" == \"Hel\"+\"lo\"", expected: true},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		testBooleanObject(t, i, evaluated, test.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		code     string
		expected float64
	}{
		{code: "1.0 + 1.0", expected: 2.0},
		{code: "1.23 * 2.0", expected: 2.46},
		{code: "1.23 * 2.0", expected: 2.46},
		{code: "(1 + 1.00) ^ 2.0", expected: 4.00},
		{code: "2.00 * 4", expected: 8.00},
		{code: "30 + 1.1", expected: 31.10},
		{code: "10 - 0.0", expected: 10.0},
		{code: "(10 + 20) ^ --1.0 / 1.0", expected: 30.0},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		testFloatObject(t, i, evaluated, test.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{code: "\"Hello\" + \" \" + \"world!\"", expected: "Hello world!"},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		testStringObject(t, i, evaluated, test.expected)
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{code: "if (10 < 20) {10} else {0}", expected: 10},
		{code: "if (10 == 20) {10} else {0}", expected: 0},
		{code: "if (10 < 20) {.10} ", expected: 0.10},
		{code: "if (31 <= 30) {0}", expected: nil},
		{code: "if (\"Hello\") {\"world!\"} ", expected: "world!"},
		{code: "if (true) {true} else {false}", expected: true},
		{code: "if (false) {true} else {false}", expected: false},
		{code: "if (1) {true} else {false}", expected: true},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		testInterface(t, i, test.expected, evaluated)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{code: "return 10;", expected: 10},
		{code: "return .8; 9 + 2;", expected: .8},
		{code: "3+2; return 10; 9 + 2;", expected: 10},
		{code: "3+2; return \"Hello\"; 9 + 2;", expected: "Hello"},
		{code: "3*20; return; 20", expected: nil},
		{code: "if (10 == 20) {10; return 0; 3;} else {1; return 20; 4;}", expected: 20},
		{code: "if (10 != 20) {if (2 > 0) {return 3;} return 1;}", expected: 3},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		testInterface(t, i, test.expected, evaluated)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{code: "\"Hello\"+3", expected: "[1,8] operator + is not defined over STRING and INTEGER"},
		{code: "\"Hello\"*.7", expected: "[1,8] operator * is not defined over STRING and FLOAT"},
		{code: "3 || 8", expected: "[1,3] operator || is not defined over INTEGERs"},
		{code: ".7 && 100.0", expected: "[1,4] operator && is not defined over FLOATs"},
		{code: "3 || 8; return 0;", expected: "[1,3] operator || is not defined over INTEGERs"},
		{code: "someVar;", expected: "[1,1] someVar is not defined"},
		{code: "len(4)", expected: "[1,4] built-in function `len` is not defined on INTEGERs"},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		err, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("case %d: no error object returned, got=%T (%+v)", i, evaluated, evaluated)
			continue
		}

		if err.Message != test.expected {
			t.Errorf("case %d: \nexpected error\t%s,\ngot\t\t%s", i, test.expected, err.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{code: "let x: Int = 5; return x;", expected: 5},
		{code: "let x: Float = .75; x;", expected: 0.75},
		{code: "let x: String = \"Hello\"; x;", expected: "Hello"},
		{code: "let _someValue: Int = (3+2)^2; _someValue;", expected: 25},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		testInterface(t, i, test.expected, evaluated)
	}
}

func TestFunctionObjects(t *testing.T) {
	input := "let addr: Func = fn (x: Int) {x + 2;};addr;"
	expectedParam := "x: Int"
	expectedBody := "(x + 2)"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not a function, got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf(
			"function has wrong Parameters, expected 1, got=%d (%+v)",
			len(fn.Parameters),
			fn.Parameters,
		)
	}

	if fn.Parameters[0].ParamString() != expectedParam {
		t.Fatalf("parameter is not %s got = %s", expectedParam, fn.Parameters[0].ParamString())
	}

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %s, got=%s", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{code: "let identity: Func = fn(x: Int) Int {return x;}; identity(5);", expected: 5},
		{code: "let double: Func = fn(x: Int) Int {return x*2;}; double(5);", expected: 10},
		{
			code:     "let add: Func = fn(x: Int, y:Int) Int {return x*2;}; add(5,add(2,3));",
			expected: 10,
		},
		{code: "fn(x:Int)Int{return x;}(5)", expected: 5},
		{code: "fn adder(x:Int) Int {return x + 10;}; adder(20);", expected: 30},
		{
			code:     "fn greeter(x: String) String {return \"Hello \" + x + \"!\";}; greeter(\"Jack\");",
			expected: "Hello Jack!",
		},
		{code: "fn isTrue(x:Bool) Bool {return x == true;}; isTrue(false);", expected: false},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		testInterface(t, i, test.expected, evaluated)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{code: `len("")`, expected: 0},
		{code: `len("four")`, expected: 4},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		testInterface(t, i, test.expected, evaluated)
	}
}

// TODO: testArrayLiteral

func TestMapLiterals(t *testing.T) {
	tests := []struct {
		code     string
		expected map[object.MapKey]int64
	}{
		{
			code: `let two: String = "two"; 
		{"one": 10 - 9, two: 2 + 0, "thr"+"ee": 6/2, 4:5, false:6}`,
			expected: map[object.MapKey]int64{
				(&object.String{Value: "one"}).MapKey():   1,
				(&object.String{Value: "two"}).MapKey():   2,
				(&object.String{Value: "three"}).MapKey(): 3,
				(&object.Integer{Value: 4}).MapKey():      5,
				FALSE.MapKey():                            6,
			},
		},
	}

	for i, test := range tests {
		evaluated := testEval(test.code)
		result, ok := evaluated.(*object.Map)
		if !ok {
			t.Fatalf("Eval didn't return Map. got=%T (%+v)", evaluated, evaluated)
		}

		if len(result.Pairs) != len(test.expected) {
			t.Fatalf("Map has wrong number of pairs, got=%d", len(result.Pairs))
		}

		for expectedKey, expectedValue := range test.expected {
			pair, ok := result.Pairs[expectedKey]
			if !ok {
				t.Errorf("No pair for given key in Pairs, %#v", expectedKey)
			}
			testIntegerObject(t, i, pair.Value, expectedValue)
		}
	}
}
