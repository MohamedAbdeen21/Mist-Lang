package eval

import (
	"lang/object"
)

func newList(elements []object.Object) *object.List {
	l := &object.List{Elements: elements}
	l.SetMethods("map", listMap)
	l.SetMethods("max", listMax)
	l.SetMethods("len", listLen)
	l.SetMethods("reverse", listReverse)
	l.SetMethods("slice", listSlice)
	l.SetMethods("filter", listFilter)
	l.SetMethods("update", listUpdate)
	return l
}

func newString(value string) *object.String {
	s := &object.String{Value: value}
	// s.SetMethods("len", stringLen)
	s.SetMethods("otherwise", stringOtherwise)
	return s
}
