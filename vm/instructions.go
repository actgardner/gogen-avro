package vm

import (
	"fmt"
)

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
