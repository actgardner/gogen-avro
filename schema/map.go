package schema

import (
	"fmt"
	"github.com/actgardner/gogen-avro/generator"
)

const mapSerializerTemplate = `
func %v(r %v, w io.Writer) error {
	err := writeLong(int64(len(r.M)), w)
	if err != nil || len(r.M) == 0 {
		return err
	}
	for k, e := range r.M {
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

const mapWrapperTemplate = `
type %[1]v struct {
	keys []string
	values []%[3]v
	M map[string]%[2]v
}

func New%[1]v() *%[1]v {
	return &%[1]v {
		keys: make([]string, 0),
		values: make([]%[3]v, 0),
		M: make(map[string]%[2]v),
	}
}

func (_ *%[1]v) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *%[1]v) SetInt(v int32) { panic("Unsupported operation") }
func (_ *%[1]v) SetLong(v int64) { panic("Unsupported operation") }
func (_ *%[1]v) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *%[1]v) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *%[1]v) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *%[1]v) SetString(v string) { panic("Unsupported operation") }
func (_ *%[1]v) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *%[1]v) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *%[1]v) SetDefault(i int) { panic("Unsupported operation") }
func (r *%[1]v) Finalize() { 
	for i := range r.keys {
		r.M[r.keys[i]] = r.values[i]
	}
	r.keys = nil
	r.values = nil
}

func (r *%[1]v) AppendMap(key string) types.Field { 
	r.keys = append(r.keys, key)
	var v %[3]v
        %[5]v
	r.values = append(r.values, v)
	return %[4]v
}

func (_ *%[1]v) AppendArray() types.Field { panic("Unsupported operation") }
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
	return fmt.Sprintf("*%v", s.Name())
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
	p.AddImport(UTIL_FILE, "github.com/actgardner/gogen-avro/vm/types")
	p.AddFunction(UTIL_FILE, s.GoType(), "", s.appendMethodDef())

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
	return fmt.Sprintf("New%v()", s.Name())
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
	return ""
}

func (s *MapField) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*MapField); ok {
		return s.ItemType().IsReadableBy(reader.ItemType())
	}
	return false
}

func (s *MapField) appendMethodDef() string {
	constructElem := ""
	ret := ""
	if constructor, ok := getConstructableForType(s.itemType); ok {
		constructElem = fmt.Sprintf("v = %v\n", constructor.ConstructorMethod())
	}
	if s.itemType.WrapperType() != "" {
		ret = fmt.Sprintf("(*%v)(&r.values[len(r.values)-1])", s.itemType.WrapperType())
	} else {
		ret = fmt.Sprintf("&r.values[len(r.values)-1]", s.GoType())
	}
	return fmt.Sprintf(mapWrapperTemplate, s.Name(), s.itemType.GoType(), s.itemType.GoType(), ret, constructElem)
}

func (s *MapField) SimpleName() string {
	return s.Name()
}
