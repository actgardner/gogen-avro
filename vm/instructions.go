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
	Call
	Return
	Halt
	Jump
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
	case Call:
		return "call"
	case Return:
		return "return"
	case Halt:
		return "halt"
	case Jump:
		return "jump"
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
	Fixed
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
	case Fixed:
		return "Fixed"
	}
	return "-"
}

const NoopField = 65535

type Instruction struct {
	Op    Op
	Type  Type
	Field int
}

func (i Instruction) String() string {
	if i.Field == NoopField {
		if i.Type.String() == "-" {
			return i.Op.String()
		}
		return fmt.Sprintf("%v(%v)", i.Op, i.Type)
	} else if i.Type.String() == "-" {
		return fmt.Sprintf("%v(%v)", i.Op, i.Field)
	}
	return fmt.Sprintf("%v(%v, %v)", i.Op, i.Type, i.Field)
}

type Program struct {
	Instructions []Instruction
}

func (p *Program) String() string {
	s := ""
	depth := ""
	for i, inst := range p.Instructions {
		if inst.Op == BlockEnd || inst.Op == SwitchEnd || inst.Op == Exit || inst.Op == SwitchCase {
			depth = depth[0 : len(depth)-3]
		}
		s += fmt.Sprintf("%v:\t%v%v\n", i, depth, inst)

		if inst.Op == BlockStart || inst.Op == SwitchStart || inst.Op == Enter || inst.Op == SwitchCase || inst.Op == AppendArray || inst.Op == AppendMap {
			depth += "|  "
		}
	}
	return s
}
