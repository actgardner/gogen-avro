package generator

import (
	"fmt"
)

const enumTypeDef = `
type %v int32

const (
%v
)

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
	return toPublicName(e.typeName) + "Enum"
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
	return fmt.Sprintf(enumTypeDef, e.GoType(), e.typeList(), e.GoType(), e.stringerList())
}

func (e *enumField) serializerMethodDef() string {
	return fmt.Sprintf(enumSerializerDef, e.SerializerMethod(), e.FieldType())
}

func (e *enumField) SerializerNs(imports, aux map[string]string) {
	aux[e.SerializerMethod()] = e.serializerMethodDef()
	aux[e.GoType()] = e.structDef()
	aux["writeInt"] = writeIntMethod
	aux["encodeInt"] = encodeIntMethod
	aux["ByteWriter"] = byteWriterInterface
}

func (e *enumField) SerializerMethod() string {
	return "write" + e.FieldType()
}
