package types

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
)

const arraySerializerTemplate = `
func %v(r %v, w io.Writer) error {
	err := writeLong(int64(len(r)),w)
	if err != nil || len(r) == 0 {
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

const arrayDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	var err error
	var blkSize int64
	var arr = make(%v, 0)
	for {
		blkSize, err = readLong(r)
		if err != nil {
			return nil, err
		}
		if blkSize == 0 {
			break
		}
		if blkSize < 0 {
			blkSize = -blkSize
			_, err = readLong(r)
			if err != nil {
				return nil, err
			}
		}
		for i := int64(0); i < blkSize; i++ {
			elem, err := %v(r)
			if err != nil {
				return nil, err
			}
			arr = append(arr, elem)
		}
	}
	return arr, nil
}
`

type arrayField struct {
	name     string
	itemType Field
}

func (s *arrayField) Name() string {
	return generator.ToPublicName(s.name)
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

func (s *arrayField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.FieldType())
}

func (s *arrayField) AddStruct(p *generator.Package) {
	s.itemType.AddStruct(p)
}

func (s *arrayField) AddSerializer(p *generator.Package) {
	itemMethodName := s.itemType.SerializerMethod()
	methodName := s.SerializerMethod()
	arraySerializer := fmt.Sprintf(arraySerializerTemplate, s.SerializerMethod(), s.GoType(), itemMethodName)
	s.itemType.AddSerializer(p)
	p.AddFunction(UTIL_FILE, "", methodName, arraySerializer)
	p.AddFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddImport(UTIL_FILE, "io")
}

func (s *arrayField) AddDeserializer(p *generator.Package) {
	itemMethodName := s.itemType.DeserializerMethod()
	methodName := s.DeserializerMethod()
	arrayDeserializer := fmt.Sprintf(arrayDeserializerTemplate, methodName, s.GoType(), s.GoType(), itemMethodName)
	s.itemType.AddDeserializer(p)
	p.AddFunction(UTIL_FILE, "", methodName, arrayDeserializer)
	p.AddFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.AddImport(UTIL_FILE, "io")
}
