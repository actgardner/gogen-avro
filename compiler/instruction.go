package compiler

import (
	"fmt"

	"github.com/actgardner/gogen-avro/vm"
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
	blockId int
}

func (b *blockStartIRInstruction) VMLength() int {
	return 8
}

// At the beginning of a block, read the length into the Long register
// If the block length is 0, jump past the block body because we're done
// If the block length is negative, read the byte count, throw it away, multiply the length by -1
// Once we've figured out the number of iterations, push the loop length onto the loop stack
func (b *blockStartIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	block := p.blocks[b.blockId]
	return []vm.Instruction{
		vm.Instruction{vm.Read, vm.Long},
		vm.Instruction{vm.EvalEqual, 0},
		vm.Instruction{vm.CondJump, block.end + 5},
		vm.Instruction{vm.EvalGreater, 0},
		vm.Instruction{vm.CondJump, block.start + 7},
		vm.Instruction{vm.Read, vm.UnusedLong},
		vm.Instruction{vm.MultLong, -1},
		vm.Instruction{vm.PushLoop, 0},
	}, nil
}

type blockEndIRInstruction struct {
	blockId int
}

func (b *blockEndIRInstruction) VMLength() int {
	return 5
}

// At the end of a block, pop the loop count and decrement it. If it's zero, go back to the very
// top to read a new block. otherwise jump to start + 7, which pushes the value back on the loop stack
func (b *blockEndIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	block := p.blocks[b.blockId]
	return []vm.Instruction{
		vm.Instruction{vm.PopLoop, 0},
		vm.Instruction{vm.AddLong, -1},
		vm.Instruction{vm.EvalEqual, 0},
		vm.Instruction{vm.CondJump, block.start},
		vm.Instruction{vm.Jump, block.start + 7},
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
	for idx, offset := range sw.caseOffsets {
		body = append(body, vm.Instruction{vm.EvalEqual, idx})
		body = append(body, vm.Instruction{vm.CondJump, offset})
	}

	body = append(body, vm.Instruction{vm.Halt, s.errId})
	return body, nil
}

type switchCaseIRInstruction struct {
	switchId    int
	writerIndex int
	// If there is no target field, or the target is not a union, the readerIndex is -1
	readerIndex int
	needsJump   bool
	insts       []vm.Instruction
}

// Creates a new switch-case IR instruction for the given switch, writer and reader IDs.
// A 'rejecting case' instruction generates a pair (Exit, Null) and a Jump past the
// switch instruction. The pair (Exit, Null) makes the VM to reject (Clear) the current
// target, which must be optional (nillable).
func newSwithCaseIRInstruction(swId, wId, rId int, isRejectingCase, needsJump bool) *switchCaseIRInstruction {
	c := &switchCaseIRInstruction{
		switchId:    swId,
		writerIndex: wId,
		readerIndex: rId,
	}
	var tmpInsts []vm.Instruction
	if isRejectingCase {
		tmpInsts = []vm.Instruction{
			vm.Instruction{vm.Jump, 0}, // 0 -> relative offset from sw.end in CompileToVM
			vm.Instruction{vm.Exit, vm.Null},
			vm.Instruction{vm.Jump, 1}, // 1 -> jump past the ending exit(Noop) in switch
		}
	} else {
		tmpInsts = []vm.Instruction{
			vm.Instruction{vm.Jump, 0},
			vm.Instruction{vm.SetLong, rId},
			vm.Instruction{vm.Set, vm.Long},
		}
	}
	fromIdx := 0
	toIdx := len(tmpInsts)
	if !needsJump {
		fromIdx++
	}
	if !isRejectingCase && rId == -1 {
		toIdx = 1
	}
	c.insts = tmpInsts[fromIdx:toIdx]
	return c
}

func (s *switchCaseIRInstruction) VMLength() int {
	return len(s.insts)
}

func (s *switchCaseIRInstruction) CompileToVM(p *irProgram) ([]vm.Instruction, error) {
	sw := p.switches[s.switchId]
	for i, _ := range s.insts {
		if s.insts[i].Op == vm.Jump {
			s.insts[i].Operand += sw.end
		}
	}
	return s.insts, nil
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
