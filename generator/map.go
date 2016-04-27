package generator

import (
	"fmt"
)

const mapSerializerTemplate = `
func %v(r %v, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = %v(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}
`

type mapField struct {
	name     string
	itemType field
}

func (s *mapField) Name() string {
	return toPublicName(s.name)
}

func (s *mapField) FieldType() string {
	return "Map" + s.itemType.FieldType()
}

func (s *mapField) GoType() string {
	return fmt.Sprintf("map[string]%v", s.itemType.GoType())
}

func (s *mapField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.FieldType())
}

func (s *mapField) AddStruct(p *Package) {}

func (s *mapField) AddSerializer(p *Package) {
	s.itemType.AddSerializer(p)
	itemMethodName := s.itemType.SerializerMethod()
	methodName := s.SerializerMethod()
	mapSerializer := fmt.Sprintf(mapSerializerTemplate, s.SerializerMethod(), s.GoType(), itemMethodName)

	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addStruct(UTIL_FILE, "StringWriter", stringWriterInterface)
	p.addFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.addFunction(UTIL_FILE, "", "writeString", writeStringMethod)
	p.addFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.addFunction(UTIL_FILE, "", methodName, mapSerializer)
	p.addImport(UTIL_FILE, "io")
}
