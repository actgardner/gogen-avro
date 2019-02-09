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
	err := program.compileType(writer, reader)
	if err != nil {
		return nil, err
	}
	return program.instructions, nil
}

func (p *Program) compileType(writer, reader schema.AvroType) error {
	//fmt.Printf("compileType(%v, %v)\n", writer, reader)
	switch writer.(type) {
	case *schema.Reference:
		if readerRef, ok := reader.(*schema.Reference); ok || reader == nil {
			return p.compileRef(writer.(*schema.Reference), readerRef)
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.MapField:
		if readerRef, ok := reader.(*schema.MapField); ok || reader == nil {
			return p.compileMap(writer.(*schema.MapField), readerRef)
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.ArrayField:
		if readerRef, ok := reader.(*schema.ArrayField); ok || reader == nil {
			return p.compileArray(writer.(*schema.ArrayField), readerRef)
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.UnionField:
		return p.compileUnion(writer.(*schema.UnionField), reader)
	case *schema.IntField:
		p.add(Read, Int, NoopField)
		if reader != nil {
			p.add(Set, Int, NoopField)
		}
		return nil
	case *schema.LongField:
		p.add(Read, Long, NoopField)
		if reader != nil {
			p.add(Set, Long, NoopField)
		}
		return nil
	case *schema.StringField:
		p.add(Read, String, NoopField)
		if reader != nil {
			p.add(Set, String, NoopField)
		}
		return nil
	case *schema.BytesField:
		p.add(Read, Bytes, NoopField)
		if reader != nil {
			p.add(Set, Bytes, NoopField)
		}
		return nil
	case *schema.FloatField:
		p.add(Read, Float, NoopField)
		if reader != nil {
			p.add(Set, Float, NoopField)
		}
		return nil
	case *schema.DoubleField:
		p.add(Read, Double, NoopField)
		if reader != nil {
			p.add(Set, Double, NoopField)
		}
		return nil
	case *schema.BoolField:
		p.add(Read, Boolean, NoopField)
		if reader != nil {
			p.add(Set, Boolean, NoopField)
		}
		return nil
	case *schema.NullField:
		return nil
	}
	return fmt.Errorf("Unsupported type: %t", writer)
}

func (p *Program) compileRef(writer, reader *schema.Reference) error {
	//fmt.Printf("compileRef(%v, %v)\n", writer, reader)
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
	case *schema.FixedDefinition:
		var readerDef *schema.FixedDefinition
		var ok bool
		if reader != nil {
			if readerDef, ok = reader.Def.(*schema.FixedDefinition); !ok {
				return fmt.Errorf("Incompatible types: %v %v", reader, writer)
			}
		}
		return p.compileFixed(writer.Def.(*schema.FixedDefinition), readerDef)
	case *schema.EnumDefinition:
		var readerDef *schema.EnumDefinition
		var ok bool
		if reader != nil {
			if readerDef, ok = reader.Def.(*schema.EnumDefinition); !ok {
				return fmt.Errorf("Incompatible types: %v %v", reader, writer)
			}
		}
		return p.compileEnum(writer.Def.(*schema.EnumDefinition), readerDef)
	}
	return fmt.Errorf("Unsupported reference type %t", reader)
}

func (p *Program) compileMap(writer, reader *schema.MapField) error {
	//fmt.Printf("compileMap(%v, %v)\n", writer, reader)
	p.add(BlockStart, Unused, NoopField)
	p.add(Read, MapKey, NoopField)
	var readerType schema.AvroType
	if reader != nil {
		p.add(AppendMap, Unused, NoopField)
		readerType = reader.ItemType()
	}
	err := p.compileType(writer.ItemType(), readerType)
	if err != nil {
		return err
	}
	if reader != nil {
		p.add(Exit, Unused, NoopField)
	}
	p.add(BlockEnd, Unused, NoopField)
	return nil
}

func (p *Program) compileArray(writer, reader *schema.ArrayField) error {
	//fmt.Printf("compileArray(%v, %v)\n", writer, reader)
	p.add(BlockStart, Unused, NoopField)
	var readerType schema.AvroType
	if reader != nil {
		p.add(AppendArray, Unused, NoopField)
		readerType = reader.ItemType()
	}
	err := p.compileType(writer.ItemType(), readerType)
	if err != nil {
		return err
	}
	if reader != nil {
		p.add(Exit, Unused, NoopField)
	}
	p.add(BlockEnd, Unused, NoopField)
	return nil
}

func (p *Program) compileRecord(writer, reader *schema.RecordDefinition) error {
	// Look up whether there's a corresonding target field and if so, parse the source field into that target
	//fmt.Printf("compileRecord(%v, %v)\n", writer, reader)
	for _, field := range writer.Fields() {
		var readerType schema.AvroType
		var readerField *schema.Field
		if reader != nil {
			readerField = reader.FieldByName(field.Name())
			if readerField != nil {
				readerType = readerField.Type()
				p.add(Enter, Unused, readerField.Index())
			}
		}
		err := p.compileType(field.Type(), readerType)
		if err != nil {
			return err
		}
		if readerField != nil {
			p.add(Exit, Unused, NoopField)
		}
	}
	return nil
}

func (p *Program) compileEnum(writer, reader *schema.EnumDefinition) error {
	//fmt.Printf("compileEnum(%v, %v)\n", writer, reader)
	p.add(Read, Int, NoopField)
	if reader != nil {
		p.add(Set, Int, NoopField)
	}
	return nil
}

func (p *Program) compileFixed(writer, reader *schema.FixedDefinition) error {
	p.add(Read, Fixed, writer.SizeBytes())
	if reader != nil {
		p.add(Set, Bytes, NoopField)
	}
	return nil
}

func (p *Program) compileUnion(writer *schema.UnionField, reader schema.AvroType) error {
	//fmt.Printf("compileUnion(%t, %t)\n", writer, reader)

	p.add(Read, UnionElem, NoopField)
	if _, ok := reader.(*schema.UnionField); ok {
		p.add(Set, UnionElem, NoopField)
	}
	p.add(SwitchStart, Unused, NoopField)
writer:
	for i, t := range writer.AvroTypes() {
		p.add(SwitchCase, Unused, i)
		if unionReader, ok := reader.(*schema.UnionField); ok {
			// If there's an exact match between the reader and writer preserve type
			// This avoids weird cases like ["string", "bytes"] which would always resolve to "string"
			if unionReader.Equals(unionReader) {
				p.add(Enter, Unused, i)
				err := p.compileType(t, writer.AvroTypes()[i])
				if err != nil {
					return err
				}
				p.add(Exit, Unused, NoopField)
				continue writer
			}
			for readerIndex, r := range unionReader.AvroTypes() {
				if t.IsReadableBy(r) {
					p.add(Enter, Unused, readerIndex)
					err := p.compileType(t, r)
					if err != nil {
						return err
					}
					p.add(Exit, Unused, NoopField)
					continue writer
				}
			}
			return fmt.Errorf("Incompatible types, no match for %v in %v", unionReader, writer)
		} else if t.IsReadableBy(reader) {
			err := p.compileType(t, reader)
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
