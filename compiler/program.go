package compiler

import (
	"github.com/actgardner/gogen-avro/vm"
)

// Build an intermediate representation of the program where
// methods, loops, switches, etc. are represented logically.
// Then concatenate everything together and replace the flow
// control with jumps to absolute offsets.

type IRProgram struct {
	main    *IRMethod
	methods map[string]*IRMethod
}

func (p *IRProgram) createMethod(name string) *IRMethod {
	method := NewIRMethod(name, p)
	p.methods[name] = method
	return method
}

func (p *IRProgram) CompileToVM() (*vm.Program, error) {
	irProgram := make([]IRInstruction, 0)

	// Concatenate all the IR instructions and assign them absolute offsets
	// Main ends with a halt, everything else ends with a ret
	p.main.addLiteral(vm.Halt, vm.Unused, vm.NoopField)
	irProgram = append(irProgram, p.main.body...)

	for _, method := range p.methods {
		method.offset = len(irProgram)
		method.addLiteral(vm.Return, vm.Unused, vm.NoopField)
		irProgram = append(irProgram, method.body...)
	}

	vmProgram := make([]vm.Instruction, 0)
	for _, instruction := range irProgram {
		compiled, err := instruction.CompileToVM(p)
		if err != nil {
			return nil, err
		}
		vmProgram = append(vmProgram, compiled)
	}
	return &vm.Program{vmProgram}, nil

}
