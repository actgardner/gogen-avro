package types

import (
	"fmt"

	"github.com/actgardner/gogen-avro/generator"
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
	itemType   AvroType
	definition map[string]interface{}
}

func NewMapField(itemType AvroType, definition map[string]interface{}) *mapField {
	return &mapField{
		itemType:   itemType,
		definition: definition,
	}
}

func (s *mapField) Name() string {
	return "Map" + s.itemType.Name()
}

func (s *mapField) SimpleName() string {
	return s.Name()
}

func (s *mapField) GoType() string {
	return fmt.Sprintf("map[string]%v", s.itemType.GoType())
}

func (s *mapField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *mapField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.Name())
}

func (s *mapField) AddStruct(p *generator.Package, containers bool) error {
	return s.itemType.AddStruct(p, containers)
}

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
	p.AddImport(UTIL_FILE, "fmt")
	p.AddImport(UTIL_FILE, "math")
}

func (s *mapField) ResolveReferences(n *Namespace) error {
	return s.itemType.ResolveReferences(n)
}

func (s *mapField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	s.definition["values"], err = s.itemType.Definition(scope)
	if err != nil {
		return nil, err
	}
	return s.definition, nil
}

func (s *mapField) ConstructorMethod() string {
	return fmt.Sprintf("make(%v)", s.GoType())
}

func (s *mapField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	items, ok := rvalue.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Expected map as default for %v, got %v", lvalue, rvalue)
	}
	setters := ""

	for k, v := range items {
		setter, err := s.itemType.DefaultValue(fmt.Sprintf("%v[%q]", lvalue, k), v)
		if err != nil {
			return "", err
		}
		setters += setter + "\n"
	}
	return setters, nil
}
