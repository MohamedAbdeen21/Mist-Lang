package eval

import (
	"lang/object"
	"math"
)

func listMax(row *int, column *int, structure object.Object, args ...object.Object) object.Object {
	list, _ := structure.(*object.List)
	elements := list.Elements
	currentType := elements[0].Type()
	maxValue := math.Inf(-1)

	for _, arg := range elements {
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

func listMap(row *int, column *int, list object.Object, args ...object.Object) object.Object {
	var newElements []object.Object
	l, _ := list.(*object.List)

	if len(args) != 1 {
		return newError("[%d,%d] map expected 1 argument, got=%d", *row, *column, len(args))
	}

	switch fn := args[0].(type) {
	case *object.Function:
		if len(fn.Parameters) != 1 {
			return newError(
				"[%d,%d] map expected its argument to have a single argument, got=%d", *row, *column, len(fn.Parameters))
		}
		for _, elem := range l.Elements {
			newElements = append(newElements, callFunction(fn, []object.Object{elem}, row, column))
		}
		return newList(newElements)
	case *object.BuiltinFunc:
		for _, elem := range l.Elements {
			newElements = append(newElements, callFunction(fn, []object.Object{elem}, row, column))
		}
		return newList(newElements)
	default:
		return newError(
			"[%d,%d] map expected its argument to be a function, got=%s",
			*row,
			*column,
			args[0].Type(),
		)
	}
}

func listFilter(row *int, column *int, list object.Object, args ...object.Object) object.Object {
	var newElements []object.Object
	l, _ := list.(*object.List)
	if len(args) != 1 {
		return newError("[%d,%d] filter expected 1 argument, got=%d", *row, *column, len(args))
	}
	switch fn := args[0].(type) {
	case *object.Function:
		if len(fn.Parameters) != 1 {
			return newError(
				"[%d,%d] filter expected its argument to have a single argument, got=%d", *row, *column, len(fn.Parameters))
		}
		if fn.ReturnType != "Bool" {
			return newError(
				"[%d,%d] filter expected its argument to return a Boolean, got=%s", *row, *column, fn.ReturnType)
		}
		for _, elem := range l.Elements {
			if callFunction(fn, []object.Object{elem}, row, column).(*object.Boolean).Value {
				newElements = append(newElements, elem)
			}
		}
		return newList(newElements)
	case *object.BuiltinFunc:
		for _, elem := range l.Elements {
			ret := callFunction(fn, []object.Object{elem}, row, column)
			val, ok := ret.(*object.Boolean)
			if !ok {
				return newError(
					"[%d,%d] filter expected its argument to return a Boolean, got=%s", *row, *column, ret.Type())
			}
			if val.Value {
				newElements = append(newElements, elem)
			}
		}
		return newList(newElements)
	default:
		return newError(
			"[%d,%d] map expected its argument to be a function, got=%s",
			*row,
			*column,
			args[0].Type(),
		)
	}
}

func listLen(row *int, column *int, list object.Object, args ...object.Object) object.Object {
	l, _ := list.(*object.List)
	if len(args) != 0 {
		return newError("[%d,%d] len expected %d arguments, got %d", *row, *column, 0, len(args))
	}
	return &object.Integer{Value: int64(len(l.Elements))}
}

func listSlice(row *int, column *int, list object.Object, args ...object.Object) object.Object {
	l, _ := list.(*object.List)
	if len(args) != 2 {
		return newError("[%d,%d] slice expected %d arguments, got %d", *row, *column, 2, len(args))
	}
	arg1, ok1 := args[0].(*object.Integer)
	arg2, ok2 := args[1].(*object.Integer)
	if !ok1 || !ok2 {
		return newError(
			"[%d,%d] slice expected arguments to be of type INTEGER, got=%s and %s",
			*row,
			*column,
			args[0].Type(),
			args[1].Type(),
		)
	}
	var newElements []object.Object
	for i := arg1.Value; i <= arg2.Value; i++ {
		newElements = append(newElements, l.Elements[i])
	}
	return newList(newElements)
}

func listReverse(row *int, column *int, list object.Object, args ...object.Object) object.Object {
	l, _ := list.(*object.List)
	if len(args) != 0 {
		return newError("[%d,%d] slice expected %d arguments, got %d", *row, *column, 0, len(args))
	}
	var newElements []object.Object
	for i := len(l.Elements) - 1; i >= 0; i-- {
		newElements = append(newElements, l.Elements[i])
	}
	return newList(newElements)
}

func listUpdate(row *int, column *int, list object.Object, args ...object.Object) object.Object {
	l, _ := list.(*object.List)
	if len(args) != 2 {
		return newError("[%d,%d] update expected %d arguments, got %d", *row, *column, 2, len(args))
	}
	arg1, ok := args[0].(*object.Integer)
	if !ok {
		return newError(
			"[%d,%d] update expected first argument to be of type INTEGER, got=%s",
			*row,
			*column,
			args[0].Type(),
		)
	}
	var newElements []object.Object
	for i := 0; i < len(l.Elements); i++ {
		if i == int(arg1.Value) {
			newElements = append(newElements, args[1])
		} else {
			newElements = append(newElements, l.Elements[i])
		}
	}
	return newList(newElements)
}
