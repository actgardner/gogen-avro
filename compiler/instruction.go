package compiler

import (
	"fmt"

	"github.com/actgardner/gogen-avro/vm"
)

type IRInstruction interface {
	CompileToVM(*IRProgram) (vm.Instruction, error)
}

type LiteralIRInstruction struct {
	instruction vm.Instruction
}

func (b *LiteralIRInstruction) CompileToVM(_ *IRProgram) (vm.Instruction, error) {
	return b.instruction, nil
}

type MethodCallIRInstruction struct {
	method string
}

func (b *MethodCallIRInstruction) CompileToVM(p *IRProgram) (vm.Instruction, error) {
	method, ok := p.methods[b.method]
	if !ok {
		return vm.Instruction{}, fmt.Errorf("Unable to call unknown method %q", b.method)
	}
	return vm.Instruction{vm.Call, vm.Unused, method.offset}, nil
}
