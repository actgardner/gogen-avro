package compiler

import (
	"fmt"

	"github.com/actgardner/gogen-avro/vm"
)

type IRInstruction interface {
	VMLength() int
	CompileToVM(*IRProgram) ([]vm.Instruction, error)
}

type LiteralIRInstruction struct {
	instruction vm.Instruction
}

func (b *LiteralIRInstruction) VMLength() int {
	return 1
}

func (b *LiteralIRInstruction) CompileToVM(_ *IRProgram) ([]vm.Instruction, error) {
	return []vm.Instruction{b.instruction}, nil
}

type MethodCallIRInstruction struct {
	method string
}

func (b *MethodCallIRInstruction) VMLength() int {
	return 1
}

func (b *MethodCallIRInstruction) CompileToVM(p *IRProgram) ([]vm.Instruction, error) {
	method, ok := p.methods[b.method]
	if !ok {
		return nil, fmt.Errorf("Unable to call unknown method %q", b.method)
	}
	return []vm.Instruction{vm.Instruction{vm.Call, vm.Unused, method.offset}}, nil
}

type BlockStartIRInstruction struct {
	blockId int
}

func (b *BlockStartIRInstruction) VMLength() int {
	return 2
}

// At the beginning of a block, read the length into the Long register
// If the block length is 0, jump past the block body because we're done
func (b *BlockStartIRInstruction) CompileToVM(p *IRProgram) ([]vm.Instruction, error) {
	block := p.blocks[b.blockId]
	return []vm.Instruction{
		vm.Instruction{vm.Read, vm.Long, vm.NoopField},
		vm.Instruction{vm.ZeroJump, vm.Unused, block.end + 3},
	}, nil
}

type BlockEndIRInstruction struct {
	blockId int
}

func (b *BlockEndIRInstruction) VMLength() int {
	return 3
}

// At the end of a block, decrement the block count. If it's zero, go back to the very
// top to read a new block. otherwise jump to start + 2, which is the beginning of the body
func (b *BlockEndIRInstruction) CompileToVM(p *IRProgram) ([]vm.Instruction, error) {
	block := p.blocks[b.blockId]
	return []vm.Instruction{
		vm.Instruction{vm.DecrLong, vm.Unused, vm.NoopField},
		vm.Instruction{vm.ZeroJump, vm.Unused, block.start},
		vm.Instruction{vm.Jump, vm.Unused, block.start + 2},
	}, nil
}
