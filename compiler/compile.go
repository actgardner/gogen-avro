package compiler

import (
	"github.com/actgardner/gogen-avro/schema"
	"github.com/actgardner/gogen-avro/vm"
)

func CompileSchemaBytes(writer, reader []byte) (*vm.Program, error) {
	readerType, err := parseSchema(reader)
	if err != nil {
		return nil, err
	}

	writerType, err := parseSchema(writer)
	if err != nil {
		return nil, err
	}

	return Compile(writerType, readerType)
}

func parseSchema(s []byte) (schema.AvroType, error) {
	ns := schema.NewNamespace(false)
	sType, err := ns.TypeForSchema(s)
	if err != nil {
		return nil, err
	}

	err = sType.ResolveReferences(ns)
	if err != nil {
		return nil, err
	}
	return sType, nil
}

func Compile(writer, reader schema.AvroType) (*vm.Program, error) {
	log("Compile()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)

	program := &IRProgram{
		methods: make(map[string]*IRMethod),
		errors:  make([]string, 0),
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
