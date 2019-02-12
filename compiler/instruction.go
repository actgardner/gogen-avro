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
	return []vm.Instruction{vm.Instruction{vm.Call, method.offset}}, nil
}

type BlockStartIRInstruction struct {
	blockId int
}

func (b *BlockStartIRInstruction) VMLength() int {
	return 9
}

// At the beginning of a block, read the length into the Long register
// If the block length is 0, jump past the block body because we're done
// If the block length is negative, read the byte count, throw it away, multiply the length by -1
func (b *BlockStartIRInstruction) CompileToVM(p *IRProgram) ([]vm.Instruction, error) {
	block := p.blocks[b.blockId]
	return []vm.Instruction{
		vm.Instruction{vm.Read, vm.Long},
		vm.Instruction{vm.LongEqual, 0},
		vm.Instruction{vm.Cond, vm.NoopField},
		vm.Instruction{vm.Jump, block.end + 5},
		vm.Instruction{vm.EvalGreater, 0},
		vm.Instruction{vm.Cond, vm.NoopField},
		vm.Instruction{vm.Jump, block.start + 9},
		vm.Instruction{vm.Read, UnusedLong},
		vm.Instruction{vm.MultLong, -1},
	}, nil
}

type BlockEndIRInstruction struct {
	blockId int
}

func (b *BlockEndIRInstruction) VMLength() int {
	return 5
}

// At the end of a block, decrement the block count. If it's zero, go back to the very
// top to read a new block. otherwise jump to start + 2, which is the beginning of the body
func (b *BlockEndIRInstruction) CompileToVM(p *IRProgram) ([]vm.Instruction, error) {
	block := p.blocks[b.blockId]
	return []vm.Instruction{
		vm.Instruction{vm.AddLong, -1},
		vm.Instruction{vm.EvalEqual, 0},
		vm.Instruction{vm.Cond, vm.NoopField},
		vm.Instruction{vm.Jump, block.start},
		vm.Instruction{vm.Jump, block.start + 4},
	}, nil
}

type SwitchStartIRInstruction struct {
	switchId int
	size     int
	errId    int
}

func (s *SwitchStartIRInstruction) VMLength() int {
	return 2*s.size + 1
}

func (s *SwitchStartIRInstruction) CompileToVM(p *IRProgram) ([]vm.Instruction, error) {
	sw := p.switches[s.switchId]
	body := []vm.Instruction{}
	for value, offset := range sw.cases {
		body = append(body, vm.Instruction{vm.CondJump, value})
		body = append(body, vm.Instruction{vm.Jump, offset + 1})
	}

	body = append(body, vm.Instruction{vm.Halt, s.errId})
	return body, nil
}

type SwitchCaseIRInstruction struct {
	switchId int
	value    int
}

func (s *SwitchCaseIRInstruction) VMLength() int {
	return 1
}

func (s *SwitchCaseIRInstruction) CompileToVM(p *IRProgram) ([]vm.Instruction, error) {
	sw := p.switches[s.switchId]
	return []vm.Instruction{
		vm.Instruction{vm.Jump, sw.end},
	}, nil
}

type SwitchEndIRInstruction struct {
	switchId int
}

func (s *SwitchEndIRInstruction) VMLength() int {
	return 0
}

func (s *SwitchEndIRInstruction) CompileToVM(p *IRProgram) ([]vm.Instruction, error) {
	return []vm.Instruction{}, nil
}
