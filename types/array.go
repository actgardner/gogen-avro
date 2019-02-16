package types

import (
	"fmt"

	"github.com/actgardner/gogen-avro/generator"
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
	itemType   AvroType
	definition map[string]interface{}
}

func NewArrayField(itemType AvroType, definition map[string]interface{}) *arrayField {
	return &arrayField{
		itemType:   itemType,
		definition: definition,
	}
}

func (s *arrayField) Name() string {
	return "Array" + s.itemType.Name()
}

func (s *arrayField) SimpleName() string {
	return s.Name()
}

func (s *arrayField) GoType() string {
	return fmt.Sprintf("[]%v", s.itemType.GoType())
}

func (s *arrayField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *arrayField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.Name())
}

func (s *arrayField) AddStruct(p *generator.Package, container bool) error {
	return s.itemType.AddStruct(p, container)
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

func (s *arrayField) ResolveReferences(n *Namespace) error {
	return s.itemType.ResolveReferences(n)
}

func (s *arrayField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	s.definition["items"], err = s.itemType.Definition(scope)
	if err != nil {
		return nil, err
	}

	return s.definition, nil
}

func (s *arrayField) ConstructorMethod() string {
	return fmt.Sprintf("make(%v, 0)", s.GoType())
}

func (s *arrayField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	items, ok := rvalue.([]interface{})
	if !ok {
		return "", fmt.Errorf("Expected array as default for %v, got %v", lvalue, rvalue)
	}

	setters := fmt.Sprintf("%v = make(%v,%v)\n", lvalue, s.GoType(), len(items))
	for i, item := range items {
		if c, ok := getConstructableForType(s.itemType); ok {
			setters += fmt.Sprintf("%v[%v] = %v\n", lvalue, i, c.ConstructorMethod())
		}

		setter, err := s.itemType.DefaultValue(fmt.Sprintf("%v[%v]", lvalue, i), item)
		if err != nil {
			return "", err
		}

		setters += setter + "\n"
	}
	return setters, nil
}
