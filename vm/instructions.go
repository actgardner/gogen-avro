package vm

import (
	"fmt"
)

type Op int

const (
	Read Op = iota
	Set
	Enter
	Exit
	BlockStart
	BlockEnd
	AppendArray
	AppendMap
	SwitchStart
	SwitchCase
	SwitchEnd
	SetDefaults
)

func (o Op) String() string {
	switch o {
	case Read:
		return "read"
	case Set:
		return "set"
	case Enter:
		return "enter"
	case Exit:
		return "exit"
	case BlockStart:
		return "block_start"
	case BlockEnd:
		return "block_end"
	case AppendArray:
		return "append_array"
	case AppendMap:
		return "append_map"
	case SwitchStart:
		return "switch_start"
	case SwitchCase:
		return "switch_case"
	case SwitchEnd:
		return "switch_end"
	case SetDefaults:
		return "set_defaults"

	}
	return "Unknown"
}

type Type int

const (
	Unused Type = iota
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
