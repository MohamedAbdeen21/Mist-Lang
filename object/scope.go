package object

type Scope struct {
	store map[string]Object
	outer *Scope
}

func NewScope() *Scope {
	return &Scope{
		store: make(map[string]Object),
		outer: nil,
	}
}

func (s *Scope) Get(name string) (Object, bool) {
	obj, ok := s.store[name]
	if !ok && s.outer != nil {
		obj, ok = s.outer.Get(name)
	}
	return obj, ok
}

func (s *Scope) Set(name string, val Object) Object {
	s.store[name] = val
	return val
}

func NewInnerScope(outer *Scope) *Scope {
	scope := NewScope()
	scope.outer = outer
	return scope
}
