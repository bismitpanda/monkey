package eval

import (
	"fmt"
	"monkey/object"
	"os"
)

var builtins = map[string]*object.Builtin{
	"len":     {Fn: bltnLen},
	"exit":    {Fn: bltnExit},
	"push":    {Fn: bltnPush},
	"last":    {Fn: bltnLast},
	"rest":    {Fn: bltnRest},
	"puts":    {Fn: bltnPuts},
	"keys":    {Fn: bltnKeys},
	"first":   {Fn: bltnFirst},
	"toInt":   {Fn: bltnToInt},
	"values":  {Fn: bltnValues},
	"toBool":  {Fn: bltnToBool},
	"locals":  {Fn: bltnLocals},
	"globals": {Fn: bltnGlobals},
}

func bltnLen(env *object.Environment, args ...object.Object) object.Object {
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

func bltnPush(env *object.Environment, args ...object.Object) object.Object {
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

func bltnFirst(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return NULL
}

func bltnLast(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		return arr.Elements[length-1]
	}

	return NULL
}

func bltnRest(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		newElements := make([]object.Object, length-1)
		copy(newElements, arr.Elements[1:length])
		return &object.Array{Elements: newElements}
	}

	return NULL
}

func bltnPuts(env *object.Environment, args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Printf("%s ", arg.Inspect())
	}
	fmt.Print("\n")
	return NULL
}

func bltnKeys(env *object.Environment, args ...object.Object) object.Object {
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

func bltnValues(env *object.Environment, args ...object.Object) object.Object {
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

func bltnExit(env *object.Environment, args ...object.Object) object.Object {
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
		fmt.Println("exit status 0")
	}

	os.Exit(exitCode)
	return NULL
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

func bltnToInt(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want = 1", len(args))
	}

	var out int

	switch arg := args[0].(type) {
	case *object.Boolean:
		if arg.Value {
			out = 1
		} else {
			out = 0
		}

	case *object.Null:
		out = 0
	default:
		return newError("invalid argument type. got %s", arg.Type())
	}

	return &object.Integer{Value: int64(out)}
}

func bltnToBool(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want = 1", len(args))
	}

	var out bool

	switch arg := args[0].(type) {
	case *object.Integer:
		if arg.Value == 0 {
			out = false
		} else {
			out = true
		}
	case *object.String:
		if arg.Value == "" {
			out = false
		} else {
			out = true
		}
	case *object.Boolean:
		out = arg.Value
	case *object.Array:
		if len(arg.Elements) == 0 {
			out = false
		} else {
			out = true
		}
	case *object.Hash:
		if len(arg.Pairs) == 0 {
			out = false
		} else {
			out = true
		}
	case *object.Null:
		out = false
	default:
		return newError("invalid argument type. got %s", arg.Type())
	}

	if out {
		return TRUE
	} else {
		return FALSE
	}
}
