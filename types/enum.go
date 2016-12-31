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

type enumField struct {
	name         string
	typeName     string
	defaultValue string
	hasDefault   bool
	symbols      []string
}

func (e *enumField) Name() string {
	return e.name
}

func (e *enumField) FieldType() string {
	return generator.ToPublicName(e.typeName)
}

func (e *enumField) GoType() string {
	return e.FieldType()
}

func (e *enumField) typeList() string {
	typeStr := ""
	for i, t := range e.symbols {
		typeStr += fmt.Sprintf("%v %v = %v\n", generator.ToPublicName(t), e.GoType(), i)
	}
	return typeStr
}

func (e *enumField) stringerList() string {
	stringerStr := ""
	for _, t := range e.symbols {
		stringerStr += fmt.Sprintf("case %v:\n return %q\n", generator.ToPublicName(t), t)
	}
	return stringerStr
}

func (e *enumField) structDef() string {
	return fmt.Sprintf(enumTypeDef, e.GoType(), e.typeList())
}

func (e *enumField) stringerDef() string {
	return fmt.Sprintf(enumTypeStringer, e.GoType(), e.stringerList())
}

func (e *enumField) serializerMethodDef() string {
	return fmt.Sprintf(enumSerializerDef, e.SerializerMethod(), e.FieldType())
}

func (e *enumField) SerializerMethod() string {
	return "write" + e.FieldType()
}

func (e *enumField) deserializerMethodDef() string {
	return fmt.Sprintf(enumDeserializerDef, e.DeserializerMethod(), e.FieldType(), e.FieldType())
}

func (e *enumField) DeserializerMethod() string {
	return "read" + e.FieldType()
}

func (e *enumField) filename() string {
	return generator.ToSnake(e.GoType()) + ".go"
}

func (e *enumField) AddStruct(p *generator.Package) {
	p.AddStruct(e.filename(), e.GoType(), e.structDef())
	p.AddFunction(e.filename(), e.GoType(), "String", e.stringerDef())
}

func (e *enumField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeInt", writeIntMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddFunction(UTIL_FILE, "", e.SerializerMethod(), e.serializerMethodDef())
	p.AddImport(UTIL_FILE, "io")
}

func (e *enumField) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readInt", readIntMethod)
	p.AddFunction(UTIL_FILE, "", e.DeserializerMethod(), e.deserializerMethodDef())
	p.AddImport(UTIL_FILE, "io")
}
