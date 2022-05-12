package object

func NewLocalEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = outer

	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, Outer: nil}
}

type Environment struct {
	store map[string]Object
	Outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.Outer != nil {
		obj, ok = e.Outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
