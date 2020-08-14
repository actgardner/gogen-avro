package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

type printFunc func()

func main() {
	printFunction := []uint16{
		0xb804, 0x0000, 0x0290, // movl $0x2000004,%eax
		0xbf01, 0x0000, 0x0090, // movl $0x1,%edi
		0x48c7, 0xc20c, 0x0, // mov $0x13, %rdx
		0x48, 0x8d35, 0x400, 0x0, // lea 0x4(%rip), %rsi
		0xf05,                  // syscall
		0xc3cc,                 // ret
		0x4865, 0x6c6c, 0x6f20, // Hello_(whitespace)
		0x576f, 0x726c, 0x6421, 0xa, // World!
	}
	executablePrintFunc, err := syscall.Mmap(
		-1,
		0,
		128,
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,
		syscall.MAP_PRIVATE|syscall.MAP_ANON)
	if err != nil {
		fmt.Printf("mmap err: %v", err)
		return
	}

	for i := range printFunction {
		executablePrintFunc[i*2] = byte(printFunction[i] >> 8)
		executablePrintFunc[i*2+1] = byte(printFunction[i])
	}

	type printFunc func()
	unsafePrintFunc := (uintptr)(unsafe.Pointer(&executablePrintFunc))
	printer := *(*printFunc)(unsafe.Pointer(&unsafePrintFunc))
	printer()
}
