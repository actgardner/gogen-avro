package vm

import (
	"fmt"
)

type Program struct {
	Instructions []Instruction
	Errors       []string
}

func (p *Program) String() string {
	s := ""
	depth := ""
	for i, inst := range p.Instructions {
		if inst.Op == SwitchEnd || inst.Op == Exit || inst.Op == SwitchCase {
			depth = depth[0 : len(depth)-3]
		}
		s += fmt.Sprintf("%v:\t%v%v\n", i, depth, inst)

		if inst.Op == SwitchStart || inst.Op == Enter || inst.Op == SwitchCase || inst.Op == AppendArray || inst.Op == AppendMap {
			depth += "|  "
		}
	}
	return s
}
