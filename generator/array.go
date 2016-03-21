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

func (s *arrayField) SerializerNs(imports, aux map[string]string) {
	s.itemType.SerializerNs(imports, aux)
	itemMethodName := s.itemType.SerializerMethod()
	methodName := s.SerializerMethod()
	if _, ok := aux[methodName]; ok {
		return
	}
	arraySerializer := fmt.Sprintf(arraySerializerTemplate, s.SerializerMethod(), s.GoType(), itemMethodName)
	aux[methodName] = arraySerializer
	aux["writeLong"] = writeLongMethod
	aux["encodeInt"] = encodeIntMethod
	aux["ByteWriter"] = byteWriterInterface
}

func (s *arrayField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.FieldType())
}
