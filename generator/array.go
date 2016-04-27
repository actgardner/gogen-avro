package generator

import (
	"fmt"
)

const arraySerializerTemplate = `
func %v(r %v, w io.Writer) error {
	err := writeLong(int64(len(r)),w)
	if err != nil {
		return err
	}
	for _, e := range r {
		err = %v(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0,w)
}
`

type arrayField struct {
	name     string
	itemType field
}

func (s *arrayField) Name() string {
	return toPublicName(s.name)
}

func (s *arrayField) FieldType() string {
	return "Array" + s.itemType.FieldType()
}

func (s *arrayField) GoType() string {
	return fmt.Sprintf("[]%v", s.itemType.GoType())
}

func (s *arrayField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.FieldType())
}

func (s *arrayField) AddStruct(p *Package) {
	s.itemType.AddStruct(p)
}

func (s *arrayField) AddSerializer(p *Package) {
	itemMethodName := s.itemType.SerializerMethod()
	methodName := s.SerializerMethod()
	arraySerializer := fmt.Sprintf(arraySerializerTemplate, s.SerializerMethod(), s.GoType(), itemMethodName)
	s.itemType.AddSerializer(p)
	p.addFunction(UTIL_FILE, "", methodName, arraySerializer)
	p.addFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.addFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addImport(UTIL_FILE, "io")
}
