package vm

import (
	"fmt"
)

type Op int

const (
	Init Op = iota
	Read
	Set
	Enter
	Exit
)

func (o Op) String() string {
	switch o {
	case Init:
		return "init"
	case Read:
		return "read"
	case Set:
		return "set"
	case Enter:
		return "enter"
	case Exit:
		return "exit"
	}
	return "Unknown"
}

type Type int

const (
	Unused = iota
	Null
	Boolean
	Int
	Long
	Float
	Double
	Bytes
	String
	MapKey
	Length
	UnionElem
)

func (t Type) String() string {
	switch t {
	case Null:
		return "Null"
	case Boolean:
		return "Boolean"
	case Int:
		return "Int"
	case Long:
		return "Long"
	case Float:
		return "Float"
	case Double:
		return "Double"
	case Bytes:
		return "Bytes"
	case String:
		return "String"
	case MapKey:
		return "MapKey"
	case Length:
		return "Length"
	case UnionElem:
		return "UnionElem"
	}
	return "Unknown"
}

const NoopField = 65535

type Instruction struct {
	Op    Op
	Type  Type
	Field int
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s(%s, %v)", i.Op, i.Type, i.Field)
}

type Assignable interface {
	// Assign a primitive field
	SetBoolean(field int, v bool)
	SetInt(field int, v int32)
	SetLong(field int, v int64)
	SetFloat(field int, v float32)
	SetDouble(field int, v float64)
	SetBytes(field int, v []byte)
	SetString(field int, v string)

	// Initialize a nested complex type
	Init(field int)

	// Get a nested complex type so we can enter it
	Get(field int) Assignable
}
