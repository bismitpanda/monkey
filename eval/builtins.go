package eval

import (
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":     object.GetBuiltinByName("len"),
	"exit":    object.GetBuiltinByName("exit"),
	"push":    object.GetBuiltinByName("push"),
	"last":    object.GetBuiltinByName("last"),
	"rest":    object.GetBuiltinByName("rest"),
	"puts":    object.GetBuiltinByName("puts"),
	"keys":    object.GetBuiltinByName("keys"),
	"first":   object.GetBuiltinByName("first"),
	"toInt":   object.GetBuiltinByName("toInt"),
	"values":  object.GetBuiltinByName("values"),
	"toBool":  object.GetBuiltinByName("toBool"),
	"locals":  {Fn: bltnLocals},
	"globals": {Fn: bltnGlobals},
}

func bltnGlobals(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError("the function globals doesnot take arguments. got = %d", len(args))
	}

	var findGlobal func(*object.Environment) *object.Environment
	findGlobal = func(env *object.Environment) *object.Environment {
		if env.Outer() != nil {
			return findGlobal(env.Outer())
		} else {
			return env
		}
	}

	globalEnv := findGlobal(env)

	pairs := make(map[object.HashKey]object.HashPair)
	for key, value := range globalEnv.Store() {
		pairKey := &object.String{Value: key}
		pairValue := value
		pairs[pairKey.HashKey()] = object.HashPair{
			Key:   pairKey,
			Value: pairValue,
		}
	}
	return &object.Hash{Pairs: pairs}
}

func bltnLocals(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError("the function globals doesnot take arguments. got = %d", len(args))
	}

	pairs := make(map[object.HashKey]object.HashPair)
	for key, value := range env.Store() {
		pairKey := &object.String{Value: key}
		pairValue := value
		pairs[pairKey.HashKey()] = object.HashPair{
			Key:   pairKey,
			Value: pairValue,
		}
	}
	return &object.Hash{Pairs: pairs}
}
