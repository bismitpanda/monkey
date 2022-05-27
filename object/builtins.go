package object

import (
	"fmt"
	"os"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{"len", &Builtin{Fn: bltnLen}},
	{"exit", &Builtin{Fn: bltnExit}},
	{"push", &Builtin{Fn: bltnPush}},
	{"last", &Builtin{Fn: bltnLast}},
	{"rest", &Builtin{Fn: bltnRest}},
	{"puts", &Builtin{Fn: bltnPuts}},
	{"keys", &Builtin{Fn: bltnKeys}},
	{"first", &Builtin{Fn: bltnFirst}},
	{"toInt", &Builtin{Fn: bltnToInt}},
	{"values", &Builtin{Fn: bltnValues}},
	{"toBool", &Builtin{Fn: bltnToBool}},
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}

	return nil
}

func bltnLen(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want = 1", len(args))
	}
	var length int64
	switch arg := args[0].(type) {
	case *String:
		length = int64(len(arg.Value))
	case *Array:
		length = int64(len(arg.Elements))
	default:
		return newError("argument to `len` not supported, got %s", arg.Type())
	}
	return &Integer{Value: length}
}

func bltnPush(env *Environment, args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got = %d, want = 2", len(args))
	}

	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)

	newElements := make([]Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &Array{Elements: newElements}

}

func bltnFirst(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want=1", len(args))
	}

	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return nil
}

func bltnLast(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want=1", len(args))
	}

	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)

	if length > 0 {
		return arr.Elements[length-1]
	}

	return nil
}

func bltnRest(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want=1", len(args))
	}

	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)

	if length > 0 {
		newElements := make([]Object, length-1)
		copy(newElements, arr.Elements[1:length])
		return &Array{Elements: newElements}
	}

	return nil
}

func bltnPuts(env *Environment, args ...Object) Object {
	for _, arg := range args {
		fmt.Printf("%s ", arg.Inspect())
	}
	fmt.Print("\n")
	return nil
}

func bltnKeys(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got = %d, want = 1", len(args))
	}
	if args[0].Type() != HASH_OBJ {
		return newError("argument to `keys` must be HASH, got %s", args[0].Type())
	}

	hash := args[0].(*Hash)

	keys := []Object{}
	for _, pair := range hash.Pairs {
		keys = append(keys, &String{Value: pair.Key.Inspect()})
	}

	return &Array{Elements: keys}
}

func bltnValues(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got = %d, want = 1", len(args))
	}
	if args[0].Type() != HASH_OBJ {
		return newError("argument to `values` must be HASH, got %s", args[0].Type())
	}

	hash := args[0].(*Hash)

	values := []Object{}
	for _, pair := range hash.Pairs {
		values = append(values, &String{Value: pair.Value.Inspect()})
	}

	return &Array{Elements: values}
}

func bltnExit(env *Environment, args ...Object) Object {
	if len(args) != 0 && len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want = 0 or 1", len(args))
	}

	exitCode := 0
	if len(args) == 1 {
		if args[0].Type() != INTEGER_OBJ {
			return newError("argument to `exit` must me INTEGER or none, got %s", args[0].Type())
		}
		exitCode = int(args[0].(*Integer).Value)
		if exitCode < 0 || exitCode > 125 {
			return newError("invalid exit code. should be within %d to %d", 0, 125)
		}
	} else {
		fmt.Println("exit status 0")
	}

	os.Exit(exitCode)
	return nil
}

func bltnToInt(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want = 1", len(args))
	}

	var out int

	switch arg := args[0].(type) {
	case *Boolean:
		if arg.Value {
			out = 1
		} else {
			out = 0
		}

	case *Null:
		out = 0
	default:
		return newError("invalid argument type. got %s", arg.Type())
	}

	return &Integer{Value: int64(out)}
}

func bltnToBool(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want = 1", len(args))
	}

	var out bool

	switch arg := args[0].(type) {
	case *Integer:
		if arg.Value == 0 {
			out = false
		} else {
			out = true
		}
	case *String:
		if arg.Value == "" {
			out = false
		} else {
			out = true
		}
	case *Boolean:
		out = arg.Value
	case *Array:
		if len(arg.Elements) == 0 {
			out = false
		} else {
			out = true
		}
	case *Hash:
		if len(arg.Pairs) == 0 {
			out = false
		} else {
			out = true
		}
	case *Null:
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
