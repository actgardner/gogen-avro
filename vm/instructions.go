package vm

import (
	"fmt"
)

// The value of NoopField as the pperand signifies the operand is unused.
const NoopField = 65535

// Constants for the data types supported by the Read and Set operations.
// If the value is > 9 it's assumed to be the length of a Fixed type.
const (
	Unused int = iota
	Null
	Boolean
	Int
	Long
	Float
	Double
	Bytes
	String
	UnionElem
	UnusedLong
)

// Represents a single VM instruction consisting of an opcode and 0 or 1 operands.
type Instruction struct {
	Op      Op
	Operand int
}

func (i Instruction) String() string {
	if i.Operand == NoopField {
		return fmt.Sprintf("%v()", i.Op)
	}
	return fmt.Sprintf("%v(%v)", i.Op, i.Operand)
}
