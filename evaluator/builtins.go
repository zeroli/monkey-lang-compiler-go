package evaluator

import (
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":  object.GetBuiltinByName("len"),
	"push": object.GetBuiltinByName("push"),
	"puts": object.GetBuiltinByName("puts"),
}
