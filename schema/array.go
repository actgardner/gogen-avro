package schema

import (
	"fmt"

	"github.com/karol-kokoszka/gogen-avro/generator"
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

const arrayWrapperTemplate = `
type %[1]v %[2]v

func (_ *%[1]v) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *%[1]v) SetInt(v int32) { panic("Unsupported operation") }
func (_ *%[1]v) SetLong(v int64) { panic("Unsupported operation") }
func (_ *%[1]v) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *%[1]v) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *%[1]v) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *%[1]v) SetString(v string) { panic("Unsupported operation") }
func (_ *%[1]v) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *%[1]v) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *%[1]v) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *%[1]v) Finalize() { }
func (_ *%[1]v) SetDefault(i int) { panic("Unsupported operation") }
func (r *%[1]v) AppendArray() types.Field {
	var v %[3]v
        %[5]v
	*r = append(*r, v)
	return %[4]v
}
`

type ArrayField struct {
	itemType   AvroType
	definition map[string]interface{}
}

func NewArrayField(itemType AvroType, definition map[string]interface{}) *ArrayField {
	return &ArrayField{
		itemType:   itemType,
		definition: definition,
	}
}

func (s *ArrayField) Name() string {
	return "Array" + s.itemType.Name()
}

func (s *ArrayField) GoType() string {
	return fmt.Sprintf("[]%v", s.itemType.GoType())
}

func (s *ArrayField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *ArrayField) AddStruct(p *generator.Package, container bool) error {
	return s.itemType.AddStruct(p, container)
}

func (s *ArrayField) ItemType() AvroType {
	return s.itemType
}

func (s *ArrayField) AddSerializer(p *generator.Package) {
	itemMethodName := s.itemType.SerializerMethod()
	methodName := s.SerializerMethod()
	arraySerializer := fmt.Sprintf(arraySerializerTemplate, s.SerializerMethod(), s.GoType(), itemMethodName)
	s.itemType.AddSerializer(p)
	p.AddFunction(UTIL_FILE, "", methodName, arraySerializer)
	p.AddFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddImport(UTIL_FILE, "github.com/karol-kokoszka/gogen-avro/vm/types")
	p.AddFunction(UTIL_FILE, s.WrapperType(), "", s.appendMethodDef())
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddImport(UTIL_FILE, "io")
}

func (s *ArrayField) ResolveReferences(n *Namespace) error {
	return s.itemType.ResolveReferences(n)
}

func (s *ArrayField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	s.definition["items"], err = s.itemType.Definition(scope)
	if err != nil {
		return nil, err
	}

	return s.definition, nil
}

func (s *ArrayField) ConstructorMethod() string {
	return fmt.Sprintf("make(%v, 0)", s.GoType())
}

func (s *ArrayField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
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

func (s *ArrayField) WrapperType() string {
	return fmt.Sprintf("%vWrapper", s.Name())
}

func (s *ArrayField) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*ArrayField); ok {
		return s.ItemType().IsReadableBy(reader.ItemType())
	}
	return false
}

func (s *ArrayField) SimpleName() string {
	return s.Name()
}

func (s *ArrayField) appendMethodDef() string {
	constructElem := ""
	ret := ""
	if constructor, ok := getConstructableForType(s.itemType); ok {
		constructElem = fmt.Sprintf("v = %v\n", constructor.ConstructorMethod())
	}
	if s.itemType.WrapperType() != "" {
		ret = fmt.Sprintf("(*%v)(&(*r)[len(*r)-1])", s.itemType.WrapperType())
	} else {
		ret = fmt.Sprintf("(*r)[len(*r)-1]")
	}
	return fmt.Sprintf(arrayWrapperTemplate, s.WrapperType(), s.GoType(), s.itemType.GoType(), ret, constructElem)
}
