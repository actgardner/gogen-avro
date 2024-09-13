package compiler

import (
	"fmt"

	"github.com/actgardner/gogen-avro/v10/vm"
)

type irInstruction interface {
	VMLength() int
	CompileToVM(*irProgram) ([]vm.Instruction, error)
}

type literalIRInstruction struct {
	instruction vm.Instruction
}

func (b *literalIRInstruction) VMLength() int {
	return 1
}

func (b *literalIRInstruction) CompileToVM(_ *irProgram) ([]vm.Instruction, error) {
	return []vm.Instruction{b.instruction}, nil
}

type methodCallIRInstruction struct {
	method string
}

func (b *methodCallIRInstruction) VMLength() int {
	return 1
}

func (b *methodCallIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	method, ok := p.methods[b.method]
	if !ok {
		return nil, fmt.Errorf("Unable to call unknown method %q", b.method)
	}
	return []vm.Instruction{vm.Instruction{vm.Call, method.offset}}, nil
}

type blockStartIRInstruction struct {
	blockId   int
	hasTarget bool
}

func (b *blockStartIRInstruction) VMLength() int {
	return blockJumpOffset(b.hasTarget) + 1
}

func blockJumpOffset(withTarget bool) int {
	if withTarget {
		return 8
	}
	return 7
}

// At the beginning of a block, read the length into the Long register
// If the block length is 0, jump past the block body because we're done
// If the block length is negative, read the byte count, throw it away, multiply the length by -1
// Once we've figured out the number of iterations, push the loop length onto the loop stack
func (b *blockStartIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	block := p.blocks[b.blockId]

	instructions := []vm.Instruction{
		{vm.Read, vm.Long},
		{vm.EvalEqual, 0},
		{vm.CondJump, block.end + 5},
		{vm.EvalGreater, 0},
		{vm.CondJump, block.start + blockJumpOffset(b.hasTarget)},
		{vm.Read, vm.UnusedLong},
		{vm.MultLong, -1},
	}

	if b.hasTarget {
		instructions = append(instructions, vm.Instruction{vm.HintSize, vm.UnusedLong})
	}

	return append(instructions, vm.Instruction{vm.PushLoop, 0}), nil
}

type blockEndIRInstruction struct {
	blockId   int
	hasTarget bool
}

func (b *blockEndIRInstruction) VMLength() int {
	return 5
}

// At the end of a block, pop the loop count and decrement it. If it's zero, go back to the very
// top to read a new block. otherwise jump to start + 7, which pushes the value back on the loop stack
func (b *blockEndIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	block := p.blocks[b.blockId]

	return []vm.Instruction{
		{vm.PopLoop, 0},
		{vm.AddLong, -1},
		{vm.EvalEqual, 0},
		{vm.CondJump, block.start},
		{vm.Jump, block.start + blockJumpOffset(b.hasTarget)},
	}, nil
}

type switchStartIRInstruction struct {
	switchId int
	size     int
	errId    int
}

func (s *switchStartIRInstruction) VMLength() int {
	return 2*s.size + 1
}

func (s *switchStartIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	sw := p.switches[s.switchId]
	body := []vm.Instruction{}
	for value, offset := range sw.cases {
		body = append(body, vm.Instruction{vm.EvalEqual, value})
		body = append(body, vm.Instruction{vm.CondJump, offset + 1})
	}

	body = append(body, vm.Instruction{vm.Halt, s.errId})
	return body, nil
}

type switchCaseIRInstruction struct {
	switchId    int
	writerIndex int
	// If there is no target field, or the target is not a union, the readerIndex is -1
	readerIndex int
}

func (s *switchCaseIRInstruction) VMLength() int {
	if s.readerIndex == -1 {
		return 1
	}
	return 1
}

func (s *switchCaseIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	sw := p.switches[s.switchId]
	if s.readerIndex == -1 {
		return []vm.Instruction{vm.Instruction{vm.Jump, sw.end}}, nil
	}

	return []vm.Instruction{
		vm.Instruction{vm.Jump, sw.end},
	}, nil
}

type switchEndIRInstruction struct {
	switchId int
}

func (s *switchEndIRInstruction) VMLength() int {
	return 0
}

func (s *switchEndIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	return []vm.Instruction{}, nil
}
