package compiler

import (
	"github.com/actgardner/gogen-avro/schema"
	"github.com/actgardner/gogen-avro/vm"
)

func Compile(writer, reader schema.AvroType) (*vm.Program, error) {
	log("Compile()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)

	program := &IRProgram{
		methods: make(map[string]*IRMethod),
	}
	program.main = NewIRMethod("main", program)

	err := program.main.compileType(writer, reader)
	if err != nil {
		return nil, err
	}

	log("%v", program)
	compiled, err := program.CompileToVM()
	log("%v", compiled)
	return compiled, err
}
