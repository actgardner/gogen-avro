package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"syscall"
	"unsafe"

	primitive "github.com/actgardner/gogen-avro/v7/test/primitive"
	"github.com/actgardner/gogen-avro/v7/vm/types"
)

type Method struct {
	data   []byte
	ops    []Op
	offset int
}

func NewMethod() (*Method, error) {
	data, err := syscall.Mmap(
		-1,
		0,
		128,
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,
		syscall.MAP_PRIVATE|syscall.MAP_ANON)
	if err != nil {
		return nil, err
	}
	return &Method{data, nil, 0}, nil
}

func (m *Method) AppendOp(op Op) {
	m.ops = append(m.ops, op)
	copy(m.prog[offset:], op.Bytes)
	m.offset += len(op.Bytes)
}

func (m *Method) AppendCall(addr uintptr) {
	// Calculate the call relative to the next instruction and encode it as 32-bit signed int
	rip := uintptr(unsafe.Pointer(&m.prog[0])) + uintptr(m.offset) + 5
	addrDiff := int32(int64(addr) - int64(rip))
	p.appendOp(CallRIP(addrDiff))
}

func (m *Method) String() string {
	for _, op := range m.ops {
		fmt.Printf("% x\t\t%v\n", op.Bytes, op.Mnemonic)
	}
}

func (m *Method) FuncPtr() uintptr {
	return uintptr(unsafe.Pointer(&p.data))
}
