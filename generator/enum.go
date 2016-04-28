package generator

import (
	"fmt"
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
	return toPublicName(e.typeName)
}

func (e *enumField) GoType() string {
	return e.FieldType()
}

func (e *enumField) typeList() string {
	typeStr := ""
	for i, t := range e.symbols {
		typeStr += fmt.Sprintf("%v %v = %v\n", toPublicName(t), e.GoType(), i)
	}
	return typeStr
}

func (e *enumField) stringerList() string {
	stringerStr := ""
	for _, t := range e.symbols {
		stringerStr += fmt.Sprintf("case %v:\n return %q\n", toPublicName(t), t)
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
	return toSnake(e.GoType()) + ".go"
}

func (e *enumField) AddStruct(p *Package) {
	p.addStruct(e.filename(), e.GoType(), e.structDef())
	p.addFunction(e.filename(), e.GoType(), "String", e.stringerDef())
}

func (e *enumField) AddSerializer(p *Package) {
	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addFunction(UTIL_FILE, "", "writeInt", writeIntMethod)
	p.addFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.addFunction(UTIL_FILE, "", e.SerializerMethod(), e.serializerMethodDef())
	p.addImport(UTIL_FILE, "io")
}

func (e *enumField) AddDeserializer(p *Package) {
	p.addFunction(UTIL_FILE, "", "readInt", readIntMethod)
	p.addFunction(UTIL_FILE, "", e.DeserializerMethod(), e.deserializerMethodDef())
	p.addImport(UTIL_FILE, "io")
}
