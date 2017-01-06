package types

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
)

const mapSerializerTemplate = `
func %v(r %v, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil || len(r) == 0 {
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

const mapDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	m := make(%v)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		if blkSize == 0 {
			break
		}
		if blkSize < 0 {
			blkSize = -blkSize
			_, err := readLong(r)
			if err != nil {
				return nil, err
			}
		}
		for i := int64(0); i < blkSize; i++ {
			key, err := readString(r)
			if err != nil {
				return nil, err
			}
			val, err := %v(r)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
	}
	return m, nil
}
`

type mapField struct {
	name     string
	itemType Field
}

func (s *mapField) Name() string {
	return generator.ToPublicName(s.name)
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

func (s *mapField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.FieldType())
}

func (s *mapField) AddStruct(p *generator.Package) {}

func (s *mapField) AddSerializer(p *generator.Package) {
	s.itemType.AddSerializer(p)
	itemMethodName := s.itemType.SerializerMethod()
	methodName := s.SerializerMethod()
	mapSerializer := fmt.Sprintf(mapSerializerTemplate, s.SerializerMethod(), s.GoType(), itemMethodName)

	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddStruct(UTIL_FILE, "StringWriter", stringWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.AddFunction(UTIL_FILE, "", "writeString", writeStringMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddFunction(UTIL_FILE, "", methodName, mapSerializer)
	p.AddImport(UTIL_FILE, "io")
}

func (s *mapField) AddDeserializer(p *generator.Package) {
	s.itemType.AddDeserializer(p)
	itemMethodName := s.itemType.DeserializerMethod()
	methodName := s.DeserializerMethod()
	mapDeserializer := fmt.Sprintf(mapDeserializerTemplate, s.DeserializerMethod(), s.GoType(), s.GoType(), itemMethodName)

	p.AddFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.AddFunction(UTIL_FILE, "", "readString", readStringMethod)
	p.AddFunction(UTIL_FILE, "", methodName, mapDeserializer)
	p.AddImport(UTIL_FILE, "io")
}
