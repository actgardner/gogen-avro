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

func (s *mapField) SerializerNs(imports, aux map[string]string) {
	s.itemType.SerializerNs(imports, aux)
	itemMethodName := s.itemType.SerializerMethod()
	methodName := s.SerializerMethod()
	if _, ok := aux[methodName]; ok {
		return
	}
	mapSerializer := fmt.Sprintf(mapSerializerTemplate, s.SerializerMethod(), s.GoType(), itemMethodName)
	aux[methodName] = mapSerializer
	aux["writeLong"] = writeLongMethod
	aux["writeString"] = writeStringMethod
	aux["encodeInt"] = encodeIntMethod
	aux["ByteWriter"] = byteWriterInterface
	aux["StringWriter"] = stringWriterInterface
}

func (s *mapField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.FieldType())
}
