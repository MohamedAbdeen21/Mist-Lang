package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"lang/ast"
	"strings"
)

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	FLOAT_OBJ    = "FLOAT"
	STRING_OBJ   = "STRING"
	NULL_OBJ     = "NULL"
	FUNCTION_OBJ = "FUNCTION"
	RETURN_OBJ   = "RETURN"
	ERROR_OBJ    = "ERROR"
	LIST_OBJ     = "LIST"
	MAP_OBJ      = "MAP"
)

func MapTypeToObject(t string) ObjectType {
	switch t {
	case "Int":
		return INTEGER_OBJ
	case "Float":
		return FLOAT_OBJ
	case "String":
		return STRING_OBJ
	case "Bool":
		return BOOLEAN_OBJ
	case "Void":
		return NULL_OBJ
	case "Func":
		return FUNCTION_OBJ
	case "List":
		return LIST_OBJ
	case "Map":
		return MAP_OBJ
	default:
		return NULL_OBJ
	}
}

type (
	ObjectType      string
	BuiltinFunction func(row *int, column *int, args ...Object) Object
	BuiltinMethod   func(row *int, column *int, structure Object, args ...Object) Object
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value   int64
	Methods map[string]BuiltinMethod
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) SetMethods(name string, method BuiltinMethod) {
	if i.Methods == nil {
		i.Methods = make(map[string]BuiltinMethod)
	}
	i.Methods[name] = method
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "" }

type String struct {
	Value   string
	Methods map[string]BuiltinMethod
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) SetMethods(name string, method BuiltinMethod) {
	if s.Methods == nil {
		s.Methods = make(map[string]BuiltinMethod)
	}
	s.Methods[name] = method
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return e.Message }

type Function struct {
	Name       *ast.Identifier
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Scope      *Scope
	ReturnType string
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.ParamString())
	}

	out.WriteString("fn" + " ")
	if f.Name != nil {
		out.WriteString(f.Name.Value)
	}
	out.WriteString("(" + strings.Join(params, ", ") + ") ")
	out.WriteString(f.ReturnType)

	return out.String()
}

type BuiltinFunc struct {
	Fn BuiltinFunction
}

func (b *BuiltinFunc) Type() ObjectType { return FUNCTION_OBJ }
func (b *BuiltinFunc) Inspect() string  { return "builtin function" }

type BuiltinMeth struct {
	Fn     BuiltinMethod
	Caller Object
}

func (b *BuiltinMeth) Type() ObjectType { return FUNCTION_OBJ }
func (b *BuiltinMeth) Inspect() string  { return "builtin method" }

type List struct {
	Elements []Object
	Methods  map[string]BuiltinMethod
}

func (l *List) Type() ObjectType { return LIST_OBJ }
func (l *List) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range l.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

func (l *List) SetMethods(name string, fn BuiltinMethod) {
	if l.Methods == nil {
		l.Methods = make(map[string]BuiltinMethod)
	}
	l.Methods[name] = fn
}

type Hashable interface {
	MapKey() MapKey
}

type MapPair struct {
	Key   Object
	Value Object
}

type Map struct {
	Pairs map[MapKey]MapPair
}

func (m *Map) Type() ObjectType { return MAP_OBJ }
func (m *Map) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type MapKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) MapKey() MapKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return MapKey{Type: b.Type(), Value: value}
}

func (i *Integer) MapKey() MapKey {
	return MapKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (f *Float) MapKey() MapKey {
	return MapKey{Type: f.Type(), Value: uint64(f.Value)}
}

func (s *String) MapKey() MapKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return MapKey{Type: s.Type(), Value: h.Sum64()}
}
