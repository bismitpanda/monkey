package eval

import (
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":     object.GetBuiltinByName("len"),
	"exit":    object.GetBuiltinByName("len"),
	"push":    object.GetBuiltinByName("len"),
	"last":    object.GetBuiltinByName("len"),
	"rest":    object.GetBuiltinByName("len"),
	"puts":    object.GetBuiltinByName("len"),
	"keys":    object.GetBuiltinByName("len"),
	"first":   object.GetBuiltinByName("len"),
	"toInt":   object.GetBuiltinByName("len"),
	"values":  object.GetBuiltinByName("len"),
	"toBool":  object.GetBuiltinByName("len"),
	"locals":  object.GetBuiltinByName("len"),
	"globals": object.GetBuiltinByName("len"),
}
