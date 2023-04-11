package eval

import (
	"fmt"
	"io/ioutil"
	"lang/object"
	"math"
	"os"
)

var Stdout *os.File

type singleton struct {
	builtins map[string]*object.BuiltinFunc
}

var single singleton

func initBuiltins() map[string]*object.BuiltinFunc {
	if single.builtins == nil {
		single.builtins = map[string]*object.BuiltinFunc{
			"len":     {Fn: lenFn},
			"max":     {Fn: maxFn},
			"print":   {Fn: printFn},
			"println": {Fn: printlnFn},
			"range":   {Fn: rangeFn},
			"string":  {Fn: convertToStringFn},
		}
	}
	return single.builtins
}

func SetupStdout() *os.File {
	if Stdout == nil {
		Stdout, _ = ioutil.TempFile("/tmp", "")
	}
	return Stdout
}

func maxFn(row *int, column *int, args ...object.Object) object.Object {
	if len(args) == 0 {
		return newError(
			"[%d,%d] max expected at least 1 argument, got=%d", *row, *column, len(args))
	}

	currentType := args[0].Type()
	maxValue := math.Inf(-1)

	for _, arg := range args {
		switch obj := arg.(type) {
		case *object.Integer:
			if obj.Type() != currentType {
				return newError(
					"[%d,%d] max expected all arguments to be of same type, found %s and %s", *row, *column, currentType, obj.Type())
			}
			if obj.Value > int64(maxValue) {
				maxValue = float64(obj.Value)
			}
		case *object.Float:
			if obj.Type() != currentType {
				return newError("[%d,%d] max expected all arguments to be of same type, found %s and %s", *row, *column, currentType, obj.Type())
			}
			if obj.Value > maxValue {
				maxValue = obj.Value
			}
		case *object.List:
			return maxFn(row, column, obj.Elements...)
		default:
			return newError(
				"[%d,%d] max expected arguments to be of type INTEGER or FLOAT, found %s", *row, *column, obj.Type())
		}
	}
	if currentType == object.INTEGER_OBJ {
		return &object.Integer{Value: int64(maxValue)}
	} else {
		return &object.Float{Value: maxValue}
	}
}

func lenFn(row *int, column *int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("[%d,%d] expected %d arguments, got %d", *row, *column, 1, len(args))
	}
	switch arg := args[0].(type) {
	case *object.List:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Map:
		return &object.Integer{Value: int64(len(arg.Pairs))}
	default:
		return newError("[%d,%d] built-in function `len` is not defined on %ss", *row, *column, arg.Type())
	}
}

func printFn(row *int, column *int, args ...object.Object) object.Object {
	for _, arg := range args {
		Stdout.WriteString(arg.Inspect())
	}
	return NULL
}

func printlnFn(row *int, column *int, args ...object.Object) object.Object {
	printFn(row, column, args...)
	Stdout.Write([]byte{'\n'})
	return NULL
}

func rangeFn(row *int, column *int, args ...object.Object) object.Object {
	if len(args) == 0 {
		return newError(
			"[%d,%d] range expected 2 arguments, got=%d", *row, *column, len(args))
	}

	arg1, ok1 := args[0].(*object.Integer)
	arg2, ok2 := args[1].(*object.Integer)
	if !ok1 || !ok2 {
		return newError(
			"[%d,%d] range expected arguments to be of type INTEGER, got=%s and %s",
			*row,
			*column,
			args[0].Type(),
			args[1].Type(),
		)
	}

	var newElements []object.Object
	for i := arg1.Value; i <= arg2.Value; i++ {
		newElements = append(newElements, &object.Integer{Value: i})
	}
	return newList(newElements)
}

func convertToStringFn(row *int, column *int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("[%d,%d] string expected %d argument, got %d", *row, *column, 1, len(args))
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		return newString(fmt.Sprintf("%d", arg.Value))
	case *object.Float:
		return newString(fmt.Sprintf("%f", arg.Value))
	case *object.Boolean:
		return newString(fmt.Sprintf("%t", arg.Value))
	default:
		return newError("[%d,%d] string can't convert value of type %s", *row, *column, arg.Type())
	}
}
