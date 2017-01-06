package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
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
	name       string
	hasDefault bool
}

func (s *nullField) Name() string {
	return generator.ToPublicName(s.name)
}

func (s *nullField) FieldType() string {
	return "Null"
}

func (s *nullField) GoType() string {
	return "interface{}"
}

func (s *nullField) SerializerMethod() string {
	return "writeNull"
}

func (s *nullField) DeserializerMethod() string {
	return "readNull"
}

func (s *nullField) AddStruct(p *generator.Package) {}

func (s *nullField) AddSerializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "writeNull", writeNullMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *nullField) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readNull", readNullMethod)
	p.AddImport(UTIL_FILE, "io")
}
