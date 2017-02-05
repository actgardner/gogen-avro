package types

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
)

const enumTypeDef = `
type %v int32

const (
%v
)
`

const enumTypeStringer = `
func (e %v) String() string {
	switch e {
%v
	}
	return "Unknown"
}
`

const enumSerializerDef = `
func %v(r %v, w io.Writer) error {
	return writeInt(int32(r), w)
}
`

const enumDeserializerDef = `
func %v(r io.Reader) (%v, error) {
	val, err := readInt(r)
	return %v(val), err
}
`

type EnumDefinition struct {
	name     QualifiedName
	aliases  []QualifiedName
	symbols  []string
	metadata map[string]interface{}
}

func (e *EnumDefinition) AvroName() QualifiedName {
	return e.name
}

func (e *EnumDefinition) Aliases() []QualifiedName {
	return e.aliases
}

func (e *EnumDefinition) FieldType() string {
	return generator.ToPublicName(e.name.Name)
}

func (e *EnumDefinition) GoType() string {
	return e.FieldType()
}

func (e *EnumDefinition) typeList() string {
	typeStr := ""
	for i, t := range e.symbols {
		typeStr += fmt.Sprintf("%v %v = %v\n", generator.ToPublicName(t), e.GoType(), i)
	}
	return typeStr
}

func (e *EnumDefinition) stringerList() string {
	stringerStr := ""
	for _, t := range e.symbols {
		stringerStr += fmt.Sprintf("case %v:\n return %q\n", generator.ToPublicName(t), t)
	}
	return stringerStr
}

func (e *EnumDefinition) structDef() string {
	return fmt.Sprintf(enumTypeDef, e.GoType(), e.typeList())
}

func (e *EnumDefinition) stringerDef() string {
	return fmt.Sprintf(enumTypeStringer, e.GoType(), e.stringerList())
}

func (e *EnumDefinition) serializerMethodDef() string {
	return fmt.Sprintf(enumSerializerDef, e.SerializerMethod(), e.FieldType())
}

func (e *EnumDefinition) SerializerMethod() string {
	return "write" + e.FieldType()
}

func (e *EnumDefinition) deserializerMethodDef() string {
	return fmt.Sprintf(enumDeserializerDef, e.DeserializerMethod(), e.FieldType(), e.FieldType())
}

func (e *EnumDefinition) DeserializerMethod() string {
	return "read" + e.FieldType()
}

func (e *EnumDefinition) filename() string {
	return generator.ToSnake(e.GoType()) + ".go"
}

func (e *EnumDefinition) AddStruct(p *generator.Package) {
	p.AddStruct(e.filename(), e.GoType(), e.structDef())
	p.AddFunction(e.filename(), e.GoType(), "String", e.stringerDef())
}

func (e *EnumDefinition) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeInt", writeIntMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddFunction(UTIL_FILE, "", e.SerializerMethod(), e.serializerMethodDef())
	p.AddImport(UTIL_FILE, "io")
}

func (e *EnumDefinition) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readInt", readIntMethod)
	p.AddFunction(UTIL_FILE, "", e.DeserializerMethod(), e.deserializerMethodDef())
	p.AddImport(UTIL_FILE, "io")
}

func (s *EnumDefinition) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *EnumDefinition) Schema(names map[QualifiedName]interface{}) interface{} {
	name := s.name.String()
	if _, ok := names[s.name]; ok {
		return name
	}
	names[s.name] = 1
	return mergeMaps(map[string]interface{}{
		"type":    "enum",
		"name":    name,
		"symbols": s.symbols,
	}, s.metadata)
}
