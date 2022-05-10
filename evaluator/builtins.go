package evaluator

import "monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {Fn: builtinLen},
}

func builtinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got = %d, want = 1", len(args))
	}
	var length int64
	switch arg := args[0].(type) {
	case *object.String:
		length = int64(len(arg.Value))
	default:
		return newError("argument to `len` not supported, got %s", arg.Type())
	}
	return &object.Integer{Value: length}
}
