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

const mapDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	m := make(%v)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Decoding block size \n", blkSize)
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

func (s *mapField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.FieldType())
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

func (s *mapField) AddDeserializer(p *Package) {
	s.itemType.AddDeserializer(p)
	itemMethodName := s.itemType.DeserializerMethod()
	methodName := s.DeserializerMethod()
	mapDeserializer := fmt.Sprintf(mapDeserializerTemplate, s.DeserializerMethod(), s.GoType(), s.GoType(), itemMethodName)

	p.addFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.addFunction(UTIL_FILE, "", "readString", readStringMethod)
	p.addFunction(UTIL_FILE, "", methodName, mapDeserializer)
	p.addImport(UTIL_FILE, "io")
	p.addImport(UTIL_FILE, "fmt")
}
