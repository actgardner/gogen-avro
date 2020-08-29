package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"syscall"
	"unsafe"
)

type program struct {
	prog   []byte
	ops    []Op
	offset int
}

func newProgram() (*program, error) {
	prog, err := syscall.Mmap(
		-1,
		0,
		128,
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,
		syscall.MAP_PRIVATE|syscall.MAP_ANON)
	if err != nil {
		return nil, err
	}
	return &program{prog, nil, 0}, nil
}

func (p *program) appendOp(op Op) {
	p.ops = append(p.ops, op)
	p.offset += len(op.Bytes)
}

// call just computes the offset for a Go method and CALLs, without setting up the stack
func (p *program) call(addr uintptr) {
	// Calculate the call relative to the next instruction and encode it as 32-bit signed int
	rip := uintptr(unsafe.Pointer(&p.prog[0])) + uintptr(p.offset) + 5
	addrDiff := int32(int64(addr) - int64(rip))
	fmt.Printf("Addr: %x, rip: %x, diff: %x\n", addr, rip, addrDiff)
	p.appendOp(CallRIP(addrDiff))
}

func (p *program) callAssigner(addr uintptr) {
	p.appendOp(SubqImm(Rsp, 0x48))
	p.appendOp(MovqSourceIDRSP(Rbp, 0x40))
	p.appendOp(LeaqSourceIDRSP(Rbp, 0x40))
	p.appendOp(MovqDestIDRSP(Rax, 0x50))
	p.appendOp(MovqDestIDRSP(Rcx, 0x58))
	p.appendOp(MovqSourceIDRSP(Rax, 0x0))
	p.appendOp(MovqSourceIDRSP(Rcx, 0x8))
	p.appendOp(MovqDestIDRSP(Rax, 0x60))
	p.appendOp(MovqDestIDRSP(Rcx, 0x68))
	p.appendOp(MovqSourceIDRSP(Rax, 0x10))
	p.appendOp(MovqSourceIDRSP(Rcx, 0x18))

	p.call(addr)

	p.appendOp(MovqDestIDRSP(Rax, 0x28))
	p.appendOp(MovqDestIDRSP(Rcx, 0x20))
	p.appendOp(MovqSourceIDRSP(Rcx, 0x30))
	p.appendOp(MovqSourceIDRSP(Rax, 0x38))
	p.appendOp(MovqSourceIDRSP(Rcx, 0x70))
	p.appendOp(MovqSourceIDRSP(Rax, 0x78))
	p.appendOp(MovqDestIDRSP(Rbp, 0x40))
	p.appendOp(AddqImm(Rsp, 0x48))

	p.appendOp(Ret())
}

func (p *program) funcPtr() unsafe.Pointer {
	offset := 0
	for _, op := range p.ops {
		fmt.Printf("% x\t\t%v\n", op.Bytes, op.Mnemonic)
		copy(p.prog[offset:], op.Bytes)
		offset += len(op.Bytes)
	}
	prog := uintptr(unsafe.Pointer(&p.prog))
	return unsafe.Pointer(&prog)
}

type Assigner func(r io.Reader, f Field) error

func callAssign(r io.Reader, f Field) error {
	return assignBool(r, f)
}

func main() {
	p, err := newProgram()
	if err != nil {
		fmt.Printf("mmap error: %v", err)
	}
	p.callAssigner(reflect.ValueOf(assignBool).Pointer())

	fn := p.funcPtr()
	exeFn := *(*Assigner)(fn)
	var val bool
	target := &Boolean{&val}
	r := bytes.NewBuffer([]byte{0x1})
	err = exeFn(r, target)
	//err := callAssign(r, target)
	fmt.Printf("Result: %v %v\n", val, err)
}
