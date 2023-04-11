package eval

import (
	"fmt"
	"lang/ast"
	"lang/object"
	"lang/token"
	"math"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, scope *object.Scope) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, scope)
	case *ast.LetStatement:
		return evalLetStatement(node, scope, &node.Token.Row, &node.Token.Column)
	case *ast.Function:
		params := node.Parameters
		body := node.Body
		function := &object.Function{Name: node.Name, Parameters: params, Scope: scope, Body: body, ReturnType: node.ReturnType()}
		scope.Set(node.Name.Value, function)
		return function
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Scope: scope, Body: body, ReturnType: node.ReturnType()}
	case *ast.CallExpression:
		function := Eval(node.Function, scope)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, scope)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return callFunction(function, args, &node.Token.Row, &node.Token.Column)
	case *ast.Identifier:
		return evalIdentifier(node, scope, &node.Token.Row, &node.Token.Column)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, scope)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return newString(node.Value)
	case *ast.Boolean:
		return evalBoolean(node.Value)
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.PrefixExpression:
		right := Eval(node.Right, scope)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, &node.Token.Row, &node.Token.Column)
	case *ast.InfixExpression:
		left := Eval(node.Left, scope)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, scope)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, &node.Token.Row, &node.Token.Column)
	case *ast.IfExpression:
		return evalIfExpression(node, scope)
	case *ast.BlockStatement:
		return evalBlockStatements(node, object.NewInnerScope(scope))
	case *ast.ListLiteral:
		elements := evalExpressions(node.Elements, scope)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return newList(elements)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, scope)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.IndexExpression:
		left := Eval(node.Left, scope)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, scope)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index, &node.Token.Row, &node.Token.Column)
	case *ast.AccessExpression:
		structure := Eval(node.Struct, scope)
		return evalAccessExpression(structure, node.Attribute, &node.Token.Row, &node.Token.Column)
	case *ast.MapLiteral:
		return evalMapLiteral(node, scope, &node.Token.Row, &node.Token.Column)
	default:
		return NULL
	}
}

func evalProgram(program *ast.Program, scope *object.Scope) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, scope)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	if result.Type() == object.RETURN_OBJ {
		returnValue := unWrapReturnValue(result)
		return returnValue
	} else {
		return result
	}
}

func evalLetStatement(
	node *ast.LetStatement,
	scope *object.Scope,
	row *int,
	column *int,
) object.Object {
	val := Eval(node.Value, scope)
	if isError(val) {
		return val
	}

	t := object.MapTypeToObject(node.Name.ReturnType())
	if t != val.Type() {
		return newError("[%d,%d] type mismatch, expected value of type %s to be of type %s",
			*row, *column, val.Type(), t)
	}

	scope.Set(node.Name.Value, val)
	return NULL
}

func evalBlockStatements(block *ast.BlockStatement, scope *object.Scope) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, scope)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalBoolean(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalPrefixExpression(
	operator string,
	right object.Object,
	row *int,
	column *int,
) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right, row, column)
	default:
		return newError(
			"[%d,%d] operator %s is not recognized as a prefix",
			*row,
			*column,
			operator,
		)
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE, NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalInfixExpression(
	operator string,
	left object.Object,
	right object.Object,
	row *int,
	column *int,
) object.Object {
	switch {
	// int and int
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, row, column)
	// float and int
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		left = &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(operator, left, right, row, column)
	case right.Type() == object.INTEGER_OBJ && left.Type() == object.FLOAT_OBJ:
		right = &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(operator, left, right, row, column)
	// boolean and boolean
	case right.Type() == object.BOOLEAN_OBJ && left.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right, row, column)
	// float and float
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right, row, column)
	// string and string
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, row, column)
		// string and int
	case left.Type() == object.STRING_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalStringIntInfixExpression(operator, left, right, row, column)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringIntInfixExpression(operator, left, right, row, column)
	// string and bool
	case left.Type() == object.STRING_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalStringBoolInfixExpression(operator, left, right, row, column)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringBoolInfixExpression(operator, left, right, row, column)
	// list and list
	case left.Type() == object.LIST_OBJ && right.Type() == object.LIST_OBJ:
		return evalListInfixExpression(operator, left, right, row, column)
	// error on other
	default:
		return newError("[%d,%d] operator %s is not defined over %s and %s",
			*row, *column, operator, left.Type(), right.Type())
	}
}

func evalIntegerInfixExpression(
	operator string,
	left object.Object,
	right object.Object,
	row *int,
	column *int,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case token.PLUS:
		return &object.Integer{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Integer{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Integer{Value: leftVal * rightVal}
	case token.SLASH:
		return &object.Integer{Value: leftVal / rightVal}
	case token.MOD:
		return &object.Integer{Value: leftVal % rightVal}
	case token.POWER:
		return &object.Integer{Value: intPow(leftVal, rightVal)}
	case token.EQ:
		return evalBoolean(leftVal == rightVal)
	case token.NE:
		return evalBoolean(leftVal != rightVal)
	case token.LT:
		return evalBoolean(leftVal < rightVal)
	case token.GT:
		return evalBoolean(leftVal > rightVal)
	case token.LE:
		return evalBoolean(leftVal <= rightVal)
	case token.GE:
		return evalBoolean(leftVal >= rightVal)
	default:
		return newError("[%d,%d] operator %s is not defined over INTEGERs",
			*row,
			*column,
			operator,
		)
	}
}

func evalFloatInfixExpression(
	operator string,
	left object.Object,
	right object.Object,
	row *int,
	column *int,
) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	switch operator {
	case token.PLUS:
		return &object.Float{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Float{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Float{Value: leftVal * rightVal}
	case token.SLASH:
		return &object.Float{Value: leftVal / rightVal}
	case token.POWER:
		return &object.Float{Value: math.Pow(leftVal, rightVal)}
	case token.EQ:
		return evalBoolean(leftVal == rightVal)
	case token.NE:
		return evalBoolean(leftVal != rightVal)
	case token.LT:
		return evalBoolean(leftVal < rightVal)
	case token.GT:
		return evalBoolean(leftVal > rightVal)
	case token.LE:
		return evalBoolean(leftVal >= rightVal)
	case token.GE:
		return evalBoolean(leftVal >= rightVal)
	default:
		return newError("[%d,%d] operator %s is not defined over FLOATs",
			*row,
			*column,
			operator,
		)
	}
}

func intPow(x int64, y int64) int64 {
	return int64(math.Pow(float64(x), float64(y)))
}

func evalMinusOperatorExpression(right object.Object, row *int, column *int) object.Object {
	switch right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -(right.(*object.Integer).Value)}
	case *object.Float:
		return &object.Float{Value: -(right.(*object.Float).Value)}
	default:
		return newError("[%d,%d] operator %s is not defined over %s", *row, *column, "-", right.Type())
	}
}

func evalBooleanInfixExpression(
	operator string,
	left object.Object,
	right object.Object,
	row *int,
	column *int,
) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value
	switch operator {
	case "||":
		return evalBoolean(leftVal || rightVal)
	case "&&":
		return evalBoolean(leftVal && rightVal)
	case "!=":
		return evalBoolean(leftVal != rightVal)
	case "==":
		return evalBoolean(leftVal == rightVal)
	default:
		return newError("[%d,%d] %s is not defined over BOOLEANs", *row, *column, operator)
	}
}

func evalStringInfixExpression(
	operator string,
	left object.Object,
	right object.Object,
	row *int,
	column *int,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "==":
		return evalBoolean(leftVal == rightVal)
	case "!=":
		return evalBoolean(leftVal != rightVal)
	case "+":
		return newString(leftVal + rightVal)
	default:
		return newError("[%d,%d] %s is not defined over STRINGs", *row, *column, operator)
	}
}

func evalListInfixExpression(
	operator string,
	left object.Object,
	right object.Object,
	row *int,
	column *int,
) object.Object {
	leftVal := left.(*object.List).Elements
	rightVal := right.(*object.List).Elements
	switch operator {
	// TODO: define more operators between Lists
	// case "==":
	// 	return evalBoolean(leftVal == rightVal)
	// case "!=":
	// 	return evalBoolean(leftVal != rightVal)
	case "+":
		return newList(append(leftVal, rightVal...))
	default:
		return newError("[%d,%d] %s is not defined over LISTs", *row, *column, operator)
	}
}

func evalStringIntInfixExpression(
	operator string,
	left object.Object,
	right object.Object,
	row *int,
	column *int,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.Integer).Value
	if rightVal < 0 {
		return newError("[%d,%d] integer can't be less than 0", *row, *column)
	}

	if operator != "*" {
		return newError(
			"[%d,%d] operator %s is not defined over %s and %s",
			*row,
			*column,
			operator,
			left.Type(),
			right.Type(),
		)
	}

	var intermediate string
	for i := rightVal; i > 0; i-- {
		intermediate += leftVal
	}

	return newString(intermediate)
}

func evalStringBoolInfixExpression(
	operator string,
	left object.Object,
	right object.Object,
	row *int,
	column *int,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.Boolean).Value
	if operator != "*" {
		return newError(
			"[%d,%d] operator %s is not defined over %s and %s",
			*row,
			*column,
			operator,
			left.Type(),
			right.Type(),
		)
	}

	if rightVal {
		return newString(leftVal)
	} else {
		return newString("")
	}
}

func evalIdentifier(
	node *ast.Identifier,
	scope *object.Scope,
	row *int,
	column *int,
) object.Object {
	if val, ok := scope.Get(node.Value); ok {
		return val
	}
	if builtin, ok := initBuiltins()[node.Value]; ok {
		return builtin
	}
	return newError("[%d,%d] %s is not defined", *row, *column, node.Value)
}

func evalIfExpression(node *ast.IfExpression, scope *object.Scope) object.Object {
	condition := Eval(node.Condition, scope)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, scope)
	} else if node.Others != nil {
		return Eval(node.Others, scope)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, scope)
	} else {
		return NULL
	}
}

func evalExpressions(exps []ast.Expression, scope *object.Scope) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, scope)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIndexExpression(
	left object.Object,
	index object.Object,
	row *int,
	column *int,
) object.Object {
	switch {
	case left.Type() == object.LIST_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalListIndexExpression(left, index, row, column)
	case left.Type() == object.MAP_OBJ:
		return evalMapIndexExpression(left, index, row, column)
	default:
		return newError(
			"[%d,%d] index operator is not defined over %ss",
			*row,
			*column,
			left.Type(),
		)
	}
}

func evalListIndexExpression(
	list object.Object,
	index object.Object,
	row *int,
	column *int,
) object.Object {
	arrayObject := list.(*object.List)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newError("[%d,%d] index %d out of range, len = %d", *row, *column, idx, max)
	}
	return arrayObject.Elements[idx]
}

func evalAccessExpression(
	exp object.Object,
	method string,
	row, column *int,
) object.Object {
	switch t := exp.(type) {
	case *object.Error:
		return exp
	case *object.List:
		if fn, ok := t.Methods[method]; !ok {
			return newError(
				"[%d,%d] type %s has no method %s",
				*row,
				*column,
				exp.Type(),
				method,
			)
		} else {
			return &object.BuiltinMeth{Fn: fn, Caller: exp}
		}
	case *object.String:
		if fn, ok := t.Methods[method]; !ok {
			return newError(
				"[%d,%d] type %s has no method %s",
				*row,
				*column,
				exp.Type(),
				method,
			)
		} else {
			return &object.BuiltinMeth{Fn: fn, Caller: exp}
		}
	default:
		return newError(
			"[%d,%d] type %s has no method %s",
			*row,
			*column,
			exp.Type(),
			method,
		)
	}
}

func callFunction(fn object.Object, args []object.Object, row *int, column *int) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		// check number of arguments
		if len(function.Parameters) != len(args) {
			return newError("[%d,%d] function %s expected %d arguments, got %d",
				*row,
				*column,
				function.Name.Value,
				len(function.Parameters),
				len(args))
		}

		// check type of each argument
		extendedScope := newFunctionScope(function, args)
		for argId, arg := range args {
			t := object.MapTypeToObject(function.Parameters[argId].ReturnType())
			if t != arg.Type() {
				return newError(
					"[%d,%d] expected argument %d (%s) to be of type %s, got %s",
					*row,
					*column,
					argId,
					function.Parameters[argId],
					t,
					arg.Type(),
				)
			}
		}

		evaluated := Eval(function.Body, extendedScope)
		expectedReturn := object.MapTypeToObject(function.ReturnType)
		returnValue := unWrapReturnValue(evaluated)
		if isError(returnValue) {
			return returnValue
		}

		if expectedReturn != returnValue.Type() {
			return newError(
				"[%d,%d] expected return to be of type %s, found %s",
				*row,
				*column,
				expectedReturn,
				returnValue.Type(),
			)
		}
		return returnValue
	case *object.BuiltinFunc:
		return function.Fn(row, column, args...)
	case *object.BuiltinMeth:
		return function.Fn(row, column, function.Caller, args...)
	default:
		return newError("[%d,%d] not a function: %s", *row, *column, fn.Type())
	}
}

func evalMapLiteral(
	node *ast.MapLiteral,
	scope *object.Scope,
	row *int,
	column *int,
) object.Object {
	pairs := make(map[object.MapKey]object.MapPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, scope)
		if isError(key) {
			return key
		}

		mapKey, ok := key.(object.Hashable)
		if !ok {
			return newError("[%d,%d] can't use %s as hash key", *row, *column, key.Type())
		}

		value := Eval(valueNode, scope)
		if isError(value) {
			return value
		}

		hashed := mapKey.MapKey()
		pairs[hashed] = object.MapPair{Key: key, Value: value}
	}
	return &object.Map{Pairs: pairs}
}

func evalMapIndexExpression(
	mapObj object.Object,
	index object.Object,
	row *int,
	column *int,
) object.Object {
	mapObject, _ := mapObj.(*object.Map)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("[%d,%d] can't use %s as hash key", *row, *column, index.Type())
	}

	pair, ok := mapObject.Pairs[key.MapKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func isTruthy(condition object.Object) bool {
	switch condition {
	case NULL, FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newFunctionScope(fn *object.Function, args []object.Object) *object.Scope {
	scope := object.NewInnerScope(fn.Scope)
	for paramIndex, param := range fn.Parameters {
		scope.Set(param.Value, args[paramIndex])
	}
	return scope
}

func unWrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
