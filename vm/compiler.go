package vm

import (
	"fmt"

	"github.com/actgardner/gogen-avro/types"
)

type Program struct {
	instructions []Instruction
}

func Compile(writer, reader types.AvroType) ([]Instruction, error) {
	program := make([]Instruction, 0)
	switch writer.(type) {
	case *types.Reference:
		if readerRef, ok := reader.(*types.Reference); ok {
			return compileRef(writer.(*types.Reference), readerRef)
		}
		return nil, fmt.Errorf("Incompatible types: %v %v", reader, writer)
	}
	return program, nil
}

func compileRef(writer, reader *types.Reference) ([]Instruction, error) {
	if writer.TypeName != reader.TypeName {
		return nil, fmt.Errorf("Incompatible types: %v %v", reader, writer)
	}

	switch writer.Def.(type) {
	case *types.RecordDefinition:
		if readerDef, ok := reader.Def.(*types.RecordDefinition); ok {
			return compileRecord(writer.Def.(*types.RecordDefinition), readerDef)
		}
		return nil, fmt.Errorf("Incompatible types: %v %v", reader, writer)
	}
	return nil, fmt.Errorf("Unsupported field %v", reader)
}

func compileRecord(writer, reader *types.RecordDefinition) ([]Instruction, error) {
	program := make([]Instruction, 0)
	// Look up whether there's a corresonding target field and if so, parse the source field into that target
	for _, field := range writer.Fields() {
		readerField := reader.FieldByName(field.Name())
		p, err := compileField(field, readerField)
		if err != nil {
			return nil, err
		}
		program = append(program, p...)
	}
	return program, nil
}

func compileField(writer, reader *types.Field) ([]Instruction, error) {
	writerType := writer.Type()
	var sourceType Type

	targetIndex := NoopField
	if reader != nil {
		targetIndex = reader.Index()
	}

	switch writerType.(type) {
	case *types.Reference:
		program := []Instruction{{Enter, Unused, targetIndex}}
		deser, err := Compile(writerType, reader.Type())
		if err != nil {
			return nil, err
		}
		program = append(program, deser...)
		program = append(program, Instruction{Exit, Unused, NoopField})
		return program, nil
	case *types.IntField:
		sourceType = Int
		break
	case *types.StringField:
		sourceType = String
		break
	default:
		return nil, fmt.Errorf("Unknown type - %v", reader)
	}

	program := []Instruction{{Read, sourceType, 0}}
	if reader != nil {
		program = append(program, Instruction{Set, sourceType, reader.Index()})
	}
	return program, nil
}
