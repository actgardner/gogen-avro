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

func (s *arrayField) GoType() string {
	return fmt.Sprintf("[]%v", s.itemType.GoType())
}

func (s *arrayField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *arrayField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.Name())
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

func (s *arrayField) ResolveReferences(n *Namespace) error {
	return s.itemType.ResolveReferences(n)
}

func (s *arrayField) Definition(scope map[QualifiedName]interface{}) interface{} {
	s.definition["items"] = s.itemType.Definition(scope)
	return s.definition
}

func (s *arrayField) ConstructorMethod() string {
	return fmt.Sprintf("make(%v, 0)", s.GoType())
}

func (s *arrayField) DefaultValue(lvalue string, rvalue interface{}) string {
	items := rvalue.([]interface{})
	setter := fmt.Sprintf("%v = make(%v,%v)\n", lvalue, s.GoType(), len(items))
	for i, item := range items {
		if c, ok := getConstructableForType(s.itemType); ok {
			setter += fmt.Sprintf("%v[%v] = %v\n", lvalue, i, c.ConstructorMethod())
		}
		setter += s.itemType.DefaultValue(fmt.Sprintf("%v[%v]", lvalue, i), item) + "\n"
	}
	return setter
}
