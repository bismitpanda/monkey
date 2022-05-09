package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (*Integer) Type() ObjectType  { return INTEGER_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (*Boolean) Type() ObjectType  { return BOOLEAN_OBJ }

type Null struct{}

func (*Null) Inspect() string  { return "null" }
func (*Null) Type() ObjectType { return NULL_OBJ }
