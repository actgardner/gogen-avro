package vm

import (
	"fmt"

	"github.com/actgardner/gogen-avro/types"
)

type Program struct {
	instructions []Instruction
}

func (p *Program) add(op Op, t Type, f int) {
	p.instructions = append(p.instructions, Instruction{op, t, f})
}

func Compile(writer, reader types.AvroType) ([]Instruction, error) {
	program := &Program{make([]Instruction, 0)}
	err := program.compileType(writer, reader, 0)
	if err != nil {
		return nil, err
	}
	return program.instructions, nil
}

func (p *Program) compileType(writer, reader types.AvroType, index int) error {
	fmt.Printf("compileType(%v, %v, %v)\n", writer, reader, index)
	switch writer.(type) {
	case *types.Reference:
		if reader != nil {
			p.add(Enter, Unused, index)
		}
		if readerRef, ok := reader.(*types.Reference); ok || reader == nil {
			err := p.compileRef(writer.(*types.Reference), readerRef)
			if err != nil {
				return err
			}
			if reader != nil {
				p.add(Exit, Unused, NoopField)
			}
			return nil
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *types.MapField:
		if reader != nil {
			p.add(Enter, Unused, index)
		}
		if readerRef, ok := reader.(*types.MapField); ok || reader == nil {
			err := p.compileMap(writer.(*types.MapField), readerRef)
			if err != nil {
				return err
			}
			if reader != nil {
				p.add(Exit, Unused, NoopField)
			}
			return nil
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *types.ArrayField:
		if reader != nil {
			p.add(Enter, Unused, index)
		}
		if readerRef, ok := reader.(*types.ArrayField); ok || reader == nil {
			err := p.compileArray(writer.(*types.ArrayField), readerRef)
			if err != nil {
				return err
			}
			if reader != nil {
				p.add(Exit, Unused, NoopField)
			}
			return nil
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *types.IntField:
		p.add(Read, Int, NoopField)
		if reader != nil {
			p.add(Set, Int, index)
		}
		return nil
	case *types.StringField:
		p.add(Read, String, NoopField)
		if reader != nil {
			p.add(Set, String, index)
		}
		return nil
	}
	return nil
}

func (p *Program) compileRef(writer, reader *types.Reference) error {
	fmt.Printf("compileRef(%v, %v)\n", writer, reader)
	if reader != nil && writer.TypeName != reader.TypeName {
		return fmt.Errorf("Incompatible types by name: %v %v", reader, writer)
	}

	switch writer.Def.(type) {
	case *types.RecordDefinition:
		var readerDef *types.RecordDefinition
		var ok bool
		if reader != nil {
			if readerDef, ok = reader.Def.(*types.RecordDefinition); !ok {
				return fmt.Errorf("Incompatible types: %v %v", reader, writer)
			}
		}
		return p.compileRecord(writer.Def.(*types.RecordDefinition), readerDef)
	}
	return fmt.Errorf("Unsupported field %v", reader)
}

func (p *Program) compileMap(writer, reader *types.MapField) error {
	fmt.Printf("compileMap(%v, %v)\n", writer, reader)
	p.add(BlockStart, Unused, NoopField)
	p.add(Read, MapKey, NoopField)
	var readerType types.AvroType
	if reader != nil {
		readerType = reader.ItemType()
	}
	err := p.compileType(writer.ItemType(), readerType, 0)
	if err != nil {
		return err
	}
	p.add(BlockEnd, Unused, NoopField)
	return nil
}

func (p *Program) compileArray(writer, reader *types.ArrayField) error {
	fmt.Printf("compileArray(%v, %v)\n", writer, reader)
	p.add(BlockStart, Unused, NoopField)
	var readerType types.AvroType
	if reader != nil {
		readerType = reader.ItemType()
	}
	err := p.compileType(writer.ItemType(), readerType, 0)
	if err != nil {
		return err
	}
	p.add(BlockEnd, Unused, NoopField)
	return nil
}

func (p *Program) compileRecord(writer, reader *types.RecordDefinition) error {
	// Look up whether there's a corresonding target field and if so, parse the source field into that target
	fmt.Printf("compileRecord(%v, %v)\n", writer, reader)
	for _, field := range writer.Fields() {
		var readerField *types.Field
		if reader != nil {
			readerField = reader.FieldByName(field.Name())
		}
		err := p.compileField(field, readerField)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Program) compileField(writer, reader *types.Field) error {
	fmt.Printf("compileField(%v, %v)\n", writer, reader)
	writerType := writer.Type()

	var readerType types.AvroType
	targetIndex := NoopField
	if reader != nil {
		targetIndex = reader.Index()
		readerType = reader.Type()
	}

	return p.compileType(writerType, readerType, targetIndex)
}
