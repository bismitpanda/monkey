package evaluator

import (
	"fmt"
	"monkey/object"
	"os"
)

var builtins = map[string]*object.Builtin{
	"len":     {Fn: builtinLen},
	"exit":    {Fn: builtinExit},
	"push":    {Fn: builtinPush},
	"puts":    {Fn: builtinPuts},
	"keys":    {Fn: builtinKeys},
	"values":  {Fn: builtinValues},
	"locals":  {Fn: builtinLocals},
	"globals": {Fn: builtinGlobals},
}

func builtinLen(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want = 1", len(args))
	}
	var length int64
	switch arg := args[0].(type) {
	case *object.String:
		length = int64(len(arg.Value))
	case *object.Array:
		length = int64(len(arg.Elements))
	default:
		return newError("argument to `len` not supported, got %s", arg.Type())
	}
	return &object.Integer{Value: length}
}

func builtinPush(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got = %d, want = 2", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	newElements := make([]object.Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}

}

func builtinPuts(env *object.Environment, args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Printf("%s ", arg.Inspect())
	}

	return NULL
}

func builtinKeys(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got = %d, want = 1", len(args))
	}
	if args[0].Type() != object.HASH_OBJ {
		return newError("argument to `keys` must be HASH, got %s", args[0].Type())
	}

	hash := args[0].(*object.Hash)

	keys := []object.Object{}
	for _, pair := range hash.Pairs {
		keys = append(keys, &object.String{Value: pair.Key.Inspect()})
	}

	return &object.Array{Elements: keys}
}

func builtinValues(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got = %d, want = 1", len(args))
	}
	if args[0].Type() != object.HASH_OBJ {
		return newError("argument to `values` must be HASH, got %s", args[0].Type())
	}

	hash := args[0].(*object.Hash)

	values := []object.Object{}
	for _, pair := range hash.Pairs {
		values = append(values, &object.String{Value: pair.Value.Inspect()})
	}

	return &object.Array{Elements: values}
}

func builtinExit(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 0 && len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want = 0 or 1", len(args))
	}

	exitCode := 0
	if len(args) == 1 {
		if args[0].Type() != object.INTEGER_OBJ {
			return newError("argument to `exit` must me INTEGER or none, got %s", args[0].Type())
		}
		exitCode = int(args[0].(*object.Integer).Value)
		if exitCode < 0 || exitCode > 125 {
			return newError("invalid exit code. should be within %d to %d", 0, 125)
		}
	} else {
		fmt.Printf("exit status 0")
	}

	os.Exit(exitCode)
	return NULL
}

func builtinGlobals(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError("the function globals doesnot take arguments. got = %d", len(args))
	}
	
	findGlobal := func(env *object.Environment) *object.Environment {
		if env.Outer != nil {
			return findGlobal(env.Outer)
		} else {
			return env
		}
	}
	
	globalEnv := findGlobal(env)
	
	pairs := make(map[object.HashKey]object.HashPair)
	for key, value := range globalEnv.Store {
		pairKey := &object.String{Value: key}
		pairValue := value
		pairs[pairKey.HashKey()] = object.HashPair{
			Key:   pairKey,
			Value: pairValue,
		}
	}
	return &object.Hash{Pairs: pairs}
}

func builtinLocals(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError("the function globals doesnot take arguments. got = %d", len(args))
	}
	
	pairs := make(map[object.HashKey]object.HashPair)
	for key, value := range env.Store {
		pairKey := &object.String{Value: key}
		pairValue := value
		pairs[pairKey.HashKey()] = object.HashPair{
			Key:   pairKey,
			Value: pairValue,
		}
	}
	return &object.Hash{Pairs: pairs}
}
