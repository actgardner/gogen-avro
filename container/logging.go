package container

import (
	"fmt"
)

var (
	LoggingEnabled = false
)

func log(f string, v ...interface{}) {
	if LoggingEnabled {
		fmt.Printf(f+"\n", v...)
	}
}
