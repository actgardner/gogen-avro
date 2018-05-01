package types

import (
	"github.com/actgardner/gogen-avro/generator"
)

const writeNullMethod = `
func writeNull(_ interface{}, _ io.Writer) error {
	return nil
}
`

const readNullMethod = `
func readNull(_ io.Reader) (interface{}, error) {
	return nil, nil
}
`

type nullField struct {
	primitiveField
}

func NewNullField(definition interface{}) *nullField {
	return &nullField{primitiveField{
		definition:         definition,
		name:               "Null",
		goType:             "interface{}",
		serializerMethod:   "writeNull",
		deserializerMethod: "readNull",
	}}
}

func (s *nullField) AddSerializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "writeNull", writeNullMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *nullField) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readNull", readNullMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *nullField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	return "", nil
}
