package vm

import (
	"fmt"
)

type Program struct {
	// The list of instructions that make up the deserializer program
	Instructions []Instruction

	// A list of errors that can be triggered by halt(x), where x is the index in this array + 1
	Errors []string
}

func (p *Program) String() string {
	s := ""
	depth := ""
	for i, inst := range p.Instructions {
		// Exit ops with Null operands are unmatched, so reducing depth
		// would cause program structure's mismatch. Better inline them.
		if inst.Op == Exit && inst.Operand != Null {
			depth = depth[0 : len(depth)-3]
		}
		s += fmt.Sprintf("%v:\t%v%v\n", i, depth, inst)

		if inst.Op == Enter || inst.Op == AppendArray || inst.Op == AppendMap {
			depth += "|  "
		}
	}

	for i, err := range p.Errors {
		s += fmt.Sprintf("Error %v:\t%v\n", i+1, err)
	}
	return s
}
