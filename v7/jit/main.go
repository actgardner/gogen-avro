package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

func PrintHello() {
	fmt.Printf("Hello World!\n")
}

type program struct {
	ops    []byte
	offset int
}

func newProgram(size int) (*program, error) {
	executableFunc, err := syscall.Mmap(
		-1,
		0,
		size,
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,
		syscall.MAP_PRIVATE|syscall.MAP_ANON)
	if err != nil {
		return nil, err
	}
	return &program{executableFunc, 0}, nil
}

func (p *program) appendBytes(ops ...byte) {
	copy(p.ops[p.offset:], ops)
	p.offset += len(ops)
}

func (p *program) call(addr uintptr) {
	// Grow the stack by 8 (subq	$0x8, %rsp)
	p.appendBytes(0x48, 0x83, 0xec, 0x08)

	// Copy BP to the stack (movq	%rbp, (%rsp))
	p.appendBytes(0x48, 0x89, 0x2c, 0x24)

	// Load new RBP (leaq	(%rsp), %rbp)
	p.appendBytes(0x48, 0x8d, 0x2c, 0x24)

	// Calculate the jump relative to the next instruction and encode it as 32-bit signed int
	rip := uintptr(unsafe.Pointer(&p.ops[0])) + uintptr(p.offset) + 5
	addrDiff := int32(int64(addr) - int64(rip))
	binary.LittleEndian.PutUint32(p.ops[p.offset+1:], *((*uint32)(unsafe.Pointer(&addrDiff))))
	p.ops[p.offset] = 0xe8

	p.offset += 5

	// movq	(%rsp), %rbp
	p.appendBytes(0x48, 0x8b, 0x2c, 0x24)

	// addq	$0x8, %rsp
	p.appendBytes(0x48, 0x83, 0xc4, 0x08)
}

func (p *program) int3() {
	p.ops[p.offset] = 0xcc
	p.offset += 1
}

func (p *program) ret() {
	p.ops[p.offset] = 0xc3
	p.ops[p.offset+1] = 0xcc
	p.offset += 2
}

func (p *program) funcPtr() unsafe.Pointer {
	prog := uintptr(unsafe.Pointer(&p.ops))
	return unsafe.Pointer(&prog)
}

func main() {
	p, err := newProgram(128)
	if err != nil {
		fmt.Printf("mmap error: %v", err)
		return
	}

	p.call(reflect.ValueOf(PrintHello).Pointer())
	//p.int3()
	p.ret()
	fn := p.funcPtr()
	exeFn := *(*func())(fn)
	exeFn()
}
