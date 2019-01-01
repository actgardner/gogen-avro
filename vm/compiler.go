package vm

import (
	"fmt"

	"github.com/actgardner/gogen-avro/schema"
)

type Program struct {
	instructions []Instruction
}

func (p *Program) add(op Op, t Type, f int) {
	p.instructions = append(p.instructions, Instruction{op, t, f})
}

func Compile(writer, reader schema.AvroType) ([]Instruction, error) {
	program := &Program{make([]Instruction, 0)}
	err := program.compileType(writer, reader, 0)
	if err != nil {
		return nil, err
	}
	return program.instructions, nil
}

func (p *Program) compileType(writer, reader schema.AvroType, index int) error {
	fmt.Printf("compileType(%v, %v, %v)\n", writer, reader, index)
	switch writer.(type) {
	case *schema.Reference:
		if reader != nil {
			p.add(Enter, Unused, index)
		}
		if readerRef, ok := reader.(*schema.Reference); ok || reader == nil {
			err := p.compileRef(writer.(*schema.Reference), readerRef)
			if err != nil {
				return err
			}
			if reader != nil {
				p.add(Exit, Unused, NoopField)
			}
			return nil
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.MapField:
		if reader != nil {
			p.add(Enter, Unused, index)
		}
		if readerRef, ok := reader.(*schema.MapField); ok || reader == nil {
			err := p.compileMap(writer.(*schema.MapField), readerRef)
			if err != nil {
				return err
			}
			if reader != nil {
				p.add(Exit, Unused, NoopField)
			}
			return nil
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.ArrayField:
		if reader != nil {
			p.add(Enter, Unused, index)
		}
		if readerRef, ok := reader.(*schema.ArrayField); ok || reader == nil {
			err := p.compileArray(writer.(*schema.ArrayField), readerRef)
			if err != nil {
				return err
			}
			if reader != nil {
				p.add(Exit, Unused, NoopField)
			}
			return nil
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.UnionField:
		err := p.compileUnion(writer.(*schema.UnionField), reader, index)
		if err != nil {
			return nil
		}
	case *schema.IntField:
		p.add(Read, Int, NoopField)
		if reader != nil {
			p.add(Set, Int, index)
		}
		return nil
	case *schema.StringField:
		p.add(Read, String, NoopField)
		if reader != nil {
			p.add(Set, String, index)
		}
		return nil
	}
	return nil
}

func (p *Program) compileRef(writer, reader *schema.Reference) error {
	fmt.Printf("compileRef(%v, %v)\n", writer, reader)
	if reader != nil && writer.TypeName != reader.TypeName {
		return fmt.Errorf("Incompatible types by name: %v %v", reader, writer)
	}

	switch writer.Def.(type) {
	case *schema.RecordDefinition:
		var readerDef *schema.RecordDefinition
		var ok bool
		if reader != nil {
			if readerDef, ok = reader.Def.(*schema.RecordDefinition); !ok {
				return fmt.Errorf("Incompatible types: %v %v", reader, writer)
			}
		}
		return p.compileRecord(writer.Def.(*schema.RecordDefinition), readerDef)
	}
	return fmt.Errorf("Unsupported field %v", reader)
}

func (p *Program) compileMap(writer, reader *schema.MapField) error {
	fmt.Printf("compileMap(%v, %v)\n", writer, reader)
	p.add(BlockStart, Unused, NoopField)
	p.add(Read, MapKey, NoopField)
	var readerType schema.AvroType
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

func (p *Program) compileArray(writer, reader *schema.ArrayField) error {
	fmt.Printf("compileArray(%v, %v)\n", writer, reader)
	p.add(BlockStart, Unused, NoopField)
	var readerType schema.AvroType
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

func (p *Program) compileRecord(writer, reader *schema.RecordDefinition) error {
	// Look up whether there's a corresonding target field and if so, parse the source field into that target
	fmt.Printf("compileRecord(%v, %v)\n", writer, reader)
	for _, field := range writer.Fields() {
		var readerField *schema.Field
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

func (p *Program) compileField(writer, reader *schema.Field) error {
	fmt.Printf("compileField(%v, %v)\n", writer, reader)
	writerType := writer.Type()

	var readerType schema.AvroType
	targetIndex := NoopField
	if reader != nil {
		targetIndex = reader.Index()
		readerType = reader.Type()
	}

	return p.compileType(writerType, readerType, targetIndex)
}

func (p *Program) compileUnion(writer *schema.UnionField, reader schema.AvroType, index int) error {
	fmt.Printf("compileUnion(%v, %v)\n", writer, reader)

	p.add(Read, UnionElem, NoopField)
	p.add(SwitchStart, Unused, NoopField)
writer:
	for i, t := range writer.AvroTypes() {
		p.add(SwitchCase, Unused, i)
		if unionReader, ok := reader.(*schema.UnionField); ok {
			for _, r := range unionReader.AvroTypes() {
				if t.IsReadableBy(r) {
					err := p.compileType(t, r, index)
					if err != nil {
						return err
					}
					continue writer
				}
			}
			return fmt.Errorf("Incompatible types: %v %v", reader, writer)
		} else if t.IsReadableBy(reader) {
			err := p.compileType(t, reader, index)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Incompatible types: %v %v", reader, writer)
		}
	}
	p.add(SwitchEnd, Unused, NoopField)
	return nil
}
