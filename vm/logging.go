package vm

import (
	"fmt"
)

var (
	CompilerLoggingEnabled = false
	VMLoggingEnabled       = false
	CompilerPrintOutput    = false
)

func compilerLog(f string, v ...interface{}) {
	if CompilerLoggingEnabled {
		fmt.Printf(f+"\n", v...)
	}
}

func runtimeLog(f string, v ...interface{}) {
	if VMLoggingEnabled {
		fmt.Printf(f+"\n", v...)
	}
}

func compileOutputLog(f string, v ...interface{}) {
	if CompilerPrintOutput {
		fmt.Printf(f+"\n", v...)
	}
}
