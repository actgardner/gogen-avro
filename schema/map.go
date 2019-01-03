package schema

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

type MapField struct {
	itemType   AvroType
	definition map[string]interface{}
}

func NewMapField(itemType AvroType, definition map[string]interface{}) *MapField {
	return &MapField{
		itemType:   itemType,
		definition: definition,
	}
}

func (s *MapField) ItemType() AvroType {
	return s.itemType
}

func (s *MapField) Name() string {
	return "Map" + s.itemType.Name()
}

func (s *MapField) GoType() string {
	return fmt.Sprintf("map[string]%v", s.itemType.GoType())
}

func (s *MapField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *MapField) AddStruct(p *generator.Package, containers bool) error {
	return s.itemType.AddStruct(p, containers)
}

func (s *MapField) AddSerializer(p *generator.Package) {
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

func (s *MapField) ResolveReferences(n *Namespace) error {
	return s.itemType.ResolveReferences(n)
}

func (s *MapField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	s.definition["values"], err = s.itemType.Definition(scope)
	if err != nil {
		return nil, err
	}
	return s.definition, nil
}

func (s *MapField) ConstructorMethod() string {
	return fmt.Sprintf("make(%v)", s.GoType())
}

func (s *MapField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
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

func (s *MapField) WrapperType() string {
	return fmt.Sprintf("%vWrapper", s.Name())
}

func (s *MapField) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*MapField); ok {
		return s.ItemType().IsReadableBy(reader.ItemType())
	}
	return false
}
