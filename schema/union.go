package schema

import (
	"fmt"
	"github.com/actgardner/gogen-avro/generator"
)

const unionSerializerTemplate = `
func %v(r %v, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType{
		%v
	}
	return fmt.Errorf("invalid value for %v")
}
`

const unionConstructorTemplate = `
func %v %v {
	return &%v{}
}
`

const unionFieldTemplate = `
func (_ %[1]v) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ %[1]v) SetInt(v int32) { panic("Unsupported operation") }
func (_ %[1]v) SetFloat(v float32) { panic("Unsupported operation") }
func (_ %[1]v) SetDouble(v float64) { panic("Unsupported operation") }
func (_ %[1]v) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ %[1]v) SetString(v string) { panic("Unsupported operation") }
func (r %[1]v) SetLong(v int64) { 
	r.UnionType = (%[2]v)(v)
}
func (r %[1]v) Get(i int) types.Field {
	switch (i) {
		%[3]v
	}
	panic("Unknown field index")
}
func (_ %[1]v) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ %[1]v) AppendArray() types.Field { panic("Unsupported operation") }
func (_ %[1]v) Finalize()  { }
`

type UnionField struct {
	name       string
	itemType   []AvroType
	definition []interface{}
}

func NewUnionField(name string, itemType []AvroType, definition []interface{}) *UnionField {
	return &UnionField{
		name:       name,
		itemType:   itemType,
		definition: definition,
	}
}

func (s *UnionField) compositeFieldName() string {
	var UnionFields = "Union"
	for _, i := range s.itemType {
		UnionFields += i.Name()
	}
	return UnionFields
}

func (s *UnionField) Name() string {
	if s.name == "" {
		return generator.ToPublicName(s.compositeFieldName())
	}
	return generator.ToPublicName(s.name)
}

func (s *UnionField) AvroTypes() []AvroType {
	return s.itemType
}

func (s *UnionField) GoType() string {
	return "*" + s.Name()
}

func (s *UnionField) unionEnumType() string {
	return fmt.Sprintf("%vTypeEnum", s.Name())
}

func (s *UnionField) unionEnumDef() string {
	var unionTypes string
	for i, item := range s.itemType {
		unionTypes += fmt.Sprintf("%v %v = %v\n", s.unionEnumType()+item.Name(), s.unionEnumType(), i)
	}
	return fmt.Sprintf("type %v int\nconst(\n%v)\n", s.unionEnumType(), unionTypes)
}

func (s *UnionField) unionTypeDef() string {
	var UnionFields string
	for _, i := range s.itemType {
		UnionFields += fmt.Sprintf("%v %v\n", i.Name(), i.GoType())
	}
	UnionFields += fmt.Sprintf("UnionType %v", s.unionEnumType())
	return fmt.Sprintf("type %v struct{\n%v\n}\n", s.Name(), UnionFields)
}

func (s *UnionField) unionSerializer() string {
	switchCase := ""
	for _, t := range s.itemType {
		switchCase += fmt.Sprintf("case %v:\nreturn %v(r.%v, w)\n", s.unionEnumType()+t.Name(), t.SerializerMethod(), t.Name())
	}
	return fmt.Sprintf(unionSerializerTemplate, s.SerializerMethod(), s.GoType(), switchCase, s.GoType())
}

func (s *UnionField) FieldsMethodDef() string {
	getBody := ""
	for i, f := range s.itemType {
		getBody += fmt.Sprintf("case %v:\n", i)
		if constructor, ok := getConstructableForType(f); ok {
			getBody += fmt.Sprintf("r.%v = %v\n", f.Name(), constructor.ConstructorMethod())
		}
		if f.WrapperType() == "" {
			getBody += fmt.Sprintf("return r.%v", f.Name())
		} else {
			getBody += fmt.Sprintf("return (*%v)(&r.%v)", f.WrapperType(), f.Name())
		}
		getBody += "\nbreak\n"
	}
	return fmt.Sprintf(unionFieldTemplate, s.GoType(), s.unionEnumType(), getBody)
}

func (s *UnionField) filename() string {
	return generator.ToSnake(s.Name()) + ".go"
}

func (s *UnionField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *UnionField) AddStruct(p *generator.Package, containers bool) error {
	p.AddStruct(s.filename(), s.unionEnumType(), s.unionEnumDef())
	p.AddStruct(s.filename(), s.Name(), s.unionTypeDef())
	p.AddFunction(s.filename(), s.GoType(), s.ConstructorMethod(), s.constructorMethodDef())
	for _, f := range s.itemType {
		err := f.AddStruct(p, containers)
		if err != nil {
			return err
		}
	}
	p.AddImport(s.filename(), "github.com/actgardner/gogen-avro/vm/types")
	p.AddFunction(s.filename(), s.GoType(), "fieldTemplate", s.FieldsMethodDef())

	return nil
}

func (s *UnionField) AddSerializer(p *generator.Package) {
	p.AddImport(UTIL_FILE, "fmt")
	p.AddFunction(UTIL_FILE, "", s.SerializerMethod(), s.unionSerializer())
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddImport(UTIL_FILE, "io")
	for _, f := range s.itemType {
		f.AddSerializer(p)
	}
}

func (s *UnionField) ResolveReferences(n *Namespace) error {
	var err error
	for _, f := range s.itemType {
		err = f.ResolveReferences(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UnionField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	for i, item := range s.itemType {
		s.definition[i], err = item.Definition(scope)
		if err != nil {
			return nil, err
		}
	}
	return s.definition, nil
}

func (s *UnionField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	lvalue = fmt.Sprintf("%v.%v", lvalue, s.itemType[0].Name())
	return s.itemType[0].DefaultValue(lvalue, rvalue)
}

func (s *UnionField) WrapperType() string {
	return ""
}

func (s *UnionField) IsReadableBy(f AvroType) bool {
	// Report if *any* writer type could be deserialized by the reader
	for _, t := range s.AvroTypes() {
		if readerUnion, ok := f.(*UnionField); ok {
			for _, rt := range readerUnion.AvroTypes() {
				if t.IsReadableBy(rt) {
					return true
				}
			}
		} else {
			if t.IsReadableBy(f) {
				return true
			}
		}
	}
	return false
}

func (s *UnionField) ConstructorMethod() string {
	return fmt.Sprintf("New%v()", s.Name())
}

func (s *UnionField) constructorMethodDef() string {
	return fmt.Sprintf(unionConstructorTemplate, s.ConstructorMethod(), s.GoType(), s.Name())
}

func (s *UnionField) Equals(reader *UnionField) bool {
	if len(reader.AvroTypes()) != len(s.AvroTypes()) {
		return false
	}

	for i, t := range s.AvroTypes() {
		readerType := reader.AvroTypes()[i]
		if writerRef, ok := t.(*Reference); ok {
			if readerRef, ok := readerType.(*Reference); ok {
				if readerRef.TypeName != writerRef.TypeName {
					return false
				}
			} else {
				return false
			}
		} else if t != readerType {
			return false
		}
	}
	return true
}
