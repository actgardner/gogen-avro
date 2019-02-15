package compiler

import (
	"fmt"

	"github.com/actgardner/gogen-avro/schema"
	"github.com/actgardner/gogen-avro/vm"
)

type irMethod struct {
	name    string
	offset  int
	body    []irInstruction
	program *irProgram
}

func newIRMethod(name string, program *irProgram) *irMethod {
	return &irMethod{
		name:    name,
		body:    make([]irInstruction, 0),
		program: program,
	}
}

func (p *irMethod) addLiteral(op vm.Op, operand int) {
	p.body = append(p.body, &literalIRInstruction{vm.Instruction{op, operand}})
}

func (p *irMethod) addMethodCall(method string) {
	p.body = append(p.body, &methodCallIRInstruction{method})
}

func (p *irMethod) addBlockStart() int {
	id := len(p.program.blocks)
	p.program.blocks = append(p.program.blocks, &irBlock{})
	p.body = append(p.body, &blockStartIRInstruction{id})
	return id
}

func (p *irMethod) addBlockEnd(id int) {
	p.body = append(p.body, &blockEndIRInstruction{id})
}

func (p *irMethod) addSwitchStart(size, errorId int) int {
	id := len(p.program.switches)
	p.program.switches = append(p.program.switches, &irSwitch{0, make(map[int]int), 0})
	p.body = append(p.body, &switchStartIRInstruction{id, size, errorId})
	return id
}

func (p *irMethod) addSwitchCase(id, writerIndex, readerIndex int) {
	p.body = append(p.body, &switchCaseIRInstruction{id, writerIndex, readerIndex})
}

func (p *irMethod) addSwitchEnd(id int) {
	p.body = append(p.body, &switchEndIRInstruction{id})
}

func (p *irMethod) addError(msg string) int {
	id := len(p.program.errors) + 1
	p.program.errors = append(p.program.errors, msg)
	return id
}

func (p *irMethod) VMLength() int {
	len := 0
	for _, inst := range p.body {
		len += inst.VMLength()
	}
	return len
}

func (p *irMethod) compileType(writer, reader schema.AvroType) error {
	log("compileType()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	switch v := writer.(type) {
	case *schema.Reference:
		if readerRef, ok := reader.(*schema.Reference); ok || reader == nil {
			return p.compileRef(v, readerRef)
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.MapField:
		if readerRef, ok := reader.(*schema.MapField); ok || reader == nil {
			return p.compileMap(v, readerRef)
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.ArrayField:
		if readerRef, ok := reader.(*schema.ArrayField); ok || reader == nil {
			return p.compileArray(v, readerRef)
		}
		return fmt.Errorf("Incompatible types: %v %v", reader, writer)
	case *schema.UnionField:
		return p.compileUnion(v, reader)
	case *schema.IntField:
		p.addLiteral(vm.Read, vm.Int)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Int)
		}
		return nil
	case *schema.LongField:
		p.addLiteral(vm.Read, vm.Long)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Long)
		}
		return nil
	case *schema.StringField:
		p.addLiteral(vm.Read, vm.String)
		if reader != nil {
			p.addLiteral(vm.Set, vm.String)
		}
		return nil
	case *schema.BytesField:
		p.addLiteral(vm.Read, vm.Bytes)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Bytes)
		}
		return nil
	case *schema.FloatField:
		p.addLiteral(vm.Read, vm.Float)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Float)
		}
		return nil
	case *schema.DoubleField:
		p.addLiteral(vm.Read, vm.Double)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Double)
		}
		return nil
	case *schema.BoolField:
		p.addLiteral(vm.Read, vm.Boolean)
		if reader != nil {
			p.addLiteral(vm.Set, vm.Boolean)
		}
		return nil
	case *schema.NullField:
		return nil
	}
	return fmt.Errorf("Unsupported type: %t", writer)
}

func (p *irMethod) compileRef(writer, reader *schema.Reference) error {
	log("compileRef()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	if reader != nil && writer.TypeName != reader.TypeName {
		return fmt.Errorf("Incompatible types by name: %v %v", reader, writer)
	}

	switch writer.Def.(type) {
	case *schema.RecordDefinition:
		var readerDef *schema.RecordDefinition
		var ok bool
		recordMethodName := fmt.Sprintf("record-r-%v", writer.Def.Name())
		if reader != nil {
			if readerDef, ok = reader.Def.(*schema.RecordDefinition); !ok {
				return fmt.Errorf("Incompatible types: %v %v", reader, writer)
			}
			recordMethodName = fmt.Sprintf("record-rw-%v", writer.Def.Name())
		}

		if _, ok := p.program.methods[recordMethodName]; !ok {
			method := p.program.createMethod(recordMethodName)
			err := method.compileRecord(writer.Def.(*schema.RecordDefinition), readerDef)
			if err != nil {
				return err
			}
		}
		p.addMethodCall(recordMethodName)
		return nil
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

func (p *irMethod) compileMap(writer, reader *schema.MapField) error {
	log("compileMap()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	blockId := p.addBlockStart()
	p.addLiteral(vm.Read, vm.String)
	var readerType schema.AvroType
	if reader != nil {
		p.addLiteral(vm.AppendMap, vm.Unused)
		readerType = reader.ItemType()
	}
	err := p.compileType(writer.ItemType(), readerType)
	if err != nil {
		return err
	}
	if reader != nil {
		p.addLiteral(vm.Exit, vm.Unused)
	}
	p.addBlockEnd(blockId)
	return nil
}

func (p *irMethod) compileArray(writer, reader *schema.ArrayField) error {
	log("compileArray()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	blockId := p.addBlockStart()
	var readerType schema.AvroType
	if reader != nil {
		p.addLiteral(vm.AppendArray, vm.Unused)
		readerType = reader.ItemType()
	}
	err := p.compileType(writer.ItemType(), readerType)
	if err != nil {
		return err
	}
	if reader != nil {
		p.addLiteral(vm.Exit, vm.Unused)
	}
	p.addBlockEnd(blockId)
	return nil
}

func (p *irMethod) compileRecord(writer, reader *schema.RecordDefinition) error {
	// Look up whether there's a corresonding target field and if so, parse the source field into that target
	log("compileRecord()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	if reader != nil {
		for _, field := range reader.Fields() {
			if writerField := writer.FieldByName(field.Name()); writerField == nil {
				if !field.HasDefault() {
					return fmt.Errorf("Incompatible schemas: field %v in reader is not present in writer and has no default value", field.Name())
				}
				p.addLiteral(vm.SetDefault, field.Index())
			}
		}
	}

	for _, field := range writer.Fields() {
		var readerType schema.AvroType
		var readerField *schema.Field
		if reader != nil {
			readerField = reader.FieldByName(field.Name())
			if readerField != nil {
				readerType = readerField.Type()
				p.addLiteral(vm.Enter, readerField.Index())
			}
		}
		err := p.compileType(field.Type(), readerType)
		if err != nil {
			return err
		}
		if readerField != nil {
			p.addLiteral(vm.Exit, vm.NoopField)
		}
	}
	return nil
}

func (p *irMethod) compileEnum(writer, reader *schema.EnumDefinition) error {
	log("compileEnum()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	p.addLiteral(vm.Read, vm.Int)
	if reader != nil {
		p.addLiteral(vm.Set, vm.Int)
	}
	return nil
}

func (p *irMethod) compileFixed(writer, reader *schema.FixedDefinition) error {
	log("compileFixed()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)
	p.addLiteral(vm.Read, 11+writer.SizeBytes())
	if reader != nil {
		p.addLiteral(vm.Set, vm.Bytes)
	}
	return nil
}

func (p *irMethod) compileUnion(writer *schema.UnionField, reader schema.AvroType) error {
	log("compileUnion()\n writer:\n %v\n---\nreader: %v\n---\n", writer, reader)

	p.addLiteral(vm.Read, vm.Long)
	errId := p.addError("Unsupported type for union")
	switchId := p.addSwitchStart(len(writer.AvroTypes()), errId)
writer:
	for i, t := range writer.AvroTypes() {
		if unionReader, ok := reader.(*schema.UnionField); ok {
			for readerIndex, r := range unionReader.AvroTypes() {
				if t.IsReadableBy(r) {
					p.addSwitchCase(switchId, i, readerIndex)
					p.addLiteral(vm.Enter, readerIndex)
					err := p.compileType(t, r)
					if err != nil {
						return err
					}
					p.addLiteral(vm.Exit, vm.NoopField)
					continue writer
				}
			}
			p.addSwitchCase(switchId, i, -1)
			typedErrId := p.addError(fmt.Sprintf("Cannot read type %v from union", t.Name()))
			p.addLiteral(vm.Halt, typedErrId)
		} else if t.IsReadableBy(reader) {
			err := p.compileType(t, reader)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Incompatible types: %v %v", reader, writer)
		}
	}
	p.addSwitchEnd(switchId)
	return nil
}
