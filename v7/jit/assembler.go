package main

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type Register byte

const (
	Eax Register = 0
	Rax Register = 0
	Ecx Register = 1
	Rcx Register = 1
	Edx Register = 2
	Rdx Register = 2
	Ebx Register = 3
	Rbx Register = 3
	Esp Register = 4
	Rsp Register = 4
	Ebp Register = 5
	Rbp Register = 5
	Esi Register = 6
	Rsi Register = 6
	Edi Register = 7
	Rdi Register = 7
)

func (r Register) String() string {
	switch r {
	case Rax:
		return "rax"
	case Rcx:
		return "rcx"
	case Rdx:
		return "rdx"
	case Rbx:
		return "rbx"
	case Rsp:
		return "rsp"
	case Rbp:
		return "rbp"
	case Rsi:
		return "rsi"
	case Rdi:
		return "rdi"
	}
	panic("Unknown register")
}

type Op struct {
	Mnemonic string
	Bytes    []byte
}

func encodeModRM(mod, reg, rm byte) byte {
	return ((mod & 3) << 6) | ((reg & 7) << 3) | (rm & 7)
}

func Ret() Op {
	return Op{
		Mnemonic: fmt.Sprintf("ret"),
		Bytes:    []byte{0xc3},
	}
}

func CallRIP(rip int32) Op {
	b := make([]byte, 5)
	b[0] = 0xe8
	binary.LittleEndian.PutUint32(b[1:], *((*uint32)(unsafe.Pointer(&rip))))
	return Op{
		Mnemonic: fmt.Sprintf("callq(%#x)", rip),
		Bytes:    b,
	}
}

func LeaqSrcIDRSP(dest Register, disp byte) Op {
	return Op{
		Mnemonic: fmt.Sprintf("leaq %#x(%%%s) %%%s", disp, Rsp, dest),
		Bytes:    []byte{0x48, 0x8d, encodeModRM(1, byte(dest), byte(Rsp)), encodeModRM(0, byte(Rsp), byte(Rsp)), disp},
	}
}

func SubqImm(r Register, v byte) Op {
	return Op{
		Mnemonic: fmt.Sprintf("subq $%#x, %%%s", v, r),
		Bytes:    []byte{0x48, 0x83, encodeModRM(3, 5, byte(r)), v},
	}
}

func AddqImm(r Register, v byte) Op {
	return Op{
		Mnemonic: fmt.Sprintf("addq $%#x, %%%s", v, r),
		Bytes:    []byte{0x48, 0x83, encodeModRM(3, 0, byte(r)), v},
	}
}

// Move quadword from  + 8-bit displacement to register
// Because SIB is computed differently for RSP
func MovqSrcIDRSP(dst Register, disp byte) Op {
	return Op{
		Mnemonic: fmt.Sprintf("movq %#x(%%%s), %%%s", disp, Rsp, dst),
		Bytes:    []byte{0x48, 0x8b, encodeModRM(1, byte(dst), byte(Rsp)), encodeModRM(0, byte(Rsp), byte(Rsp)), disp},
	}
}

// Move quadword to RSP + 8-bit displacement from register
// Because SIB is computed differently for RSP
func MovqDestIDRSP(src Register, disp byte) Op {
	return Op{
		Mnemonic: fmt.Sprintf("movq %%%s, %#x(%%%s)", src, disp, Rsp),
		Bytes:    []byte{0x48, 0x89, encodeModRM(1, byte(src), byte(Rsp)), encodeModRM(0, byte(Rsp), byte(Rsp)), disp},
	}
}

/*
// Move immediate quadword to RSP
func MovqDestRSPImm(src Register, imm int32) Op {
	return Op{
		Mnemonic: fmt.Sprintf("movq %#X, (%%%s)", imm, Rsp),
		Bytes:    []byte{0xc7, 0x89, encodeModRM(1, byte(src), byte(Rsp)), encodeModRM(0, byte(Rsp), byte(Rsp)), disp},
	}
}
*/
