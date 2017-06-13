package types

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
	"io"
)

const unionSerializerTemplate = `
func %v(r %v, w io.Writer) error {
	err := types.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType{
		%v
	}
	return fmt.Errorf("Invalid value for %v")
}
`

const unionDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	field, err := types.ReadLong(r)
	var unionStr %v
	if err != nil {
		return unionStr, err
	}
	unionStr.UnionType = %v(field)
	switch unionStr.UnionType {
		%v
	default:	
		return unionStr, fmt.Errorf("Invalid value for %v")
	}
	return unionStr, nil
}
`

type unionField struct {
	itemType   []AvroType
	definition []interface{}
}

func NewUnionField(itemType []AvroType, definition []interface{}) *unionField {
	return &unionField{
		itemType:   itemType,
		definition: definition,
	}
}

func (s *unionField) Name() string {
	return s.GoType()
}

func (s *unionField) GoType() string {
	var unionFields = "Union"
	for _, i := range s.itemType {
		unionFields += i.Name()
	}
	return generator.ToPublicName(unionFields)
}

func (s *unionField) unionEnumType() string {
	return fmt.Sprintf("%vTypeEnum", s.Name())
}

func (s *unionField) unionEnumDef() string {
	var unionTypes string
	for i, item := range s.itemType {
		unionTypes += fmt.Sprintf("%v %v = %v\n", s.unionEnumType()+item.Name(), s.unionEnumType(), i)
	}
	return fmt.Sprintf("type %v int\nconst(\n%v)\n", s.unionEnumType(), unionTypes)
}

func (s *unionField) unionTypeDef() string {
	var unionFields string
	for _, i := range s.itemType {
		unionFields += fmt.Sprintf("%v %v\n", i.Name(), i.GoType())
	}
	unionFields += fmt.Sprintf("UnionType %v", s.unionEnumType())
	return fmt.Sprintf("type %v struct{\n%v\n}\n", s.Name(), unionFields)
}

func (s *unionField) unionSerializer() string {
	switchCase := ""
	for _, t := range s.itemType {
		switchCase += fmt.Sprintf("case %v:\nreturn %v(r.%v, w)\n", s.unionEnumType()+t.Name(), t.SerializerMethod(), t.Name())
	}
	return fmt.Sprintf(unionSerializerTemplate, s.SerializerMethod(), s.GoType(), switchCase, s.GoType())
}

func (s *unionField) unionDeserializer() string {
	switchCase := ""
	for _, t := range s.itemType {
		switchCase += fmt.Sprintf("case %v:\nval, err :=  %v(r)\nif err != nil {return unionStr, err}\nunionStr.%v = val\n", s.unionEnumType()+t.Name(), t.DeserializerMethod(), t.Name())
	}
	return fmt.Sprintf(unionDeserializerTemplate, s.DeserializerMethod(), s.GoType(), s.GoType(), s.unionEnumType(), switchCase, s.GoType())
}

func (s *unionField) filename() string {
	return generator.ToSnake(s.GoType()) + ".go"
}

func (s *unionField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *unionField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.Name())
}

func (s *unionField) AddStruct(p *generator.Package, containers bool) error {
	p.AddStruct(s.filename(), s.unionEnumType(), s.unionEnumDef())
	p.AddStruct(s.filename(), s.Name(), s.unionTypeDef())
	for _, f := range s.itemType {
		err := f.AddStruct(p, containers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *unionField) AddSerializer(p *generator.Package) {
	p.AddImport(UTIL_FILE, "fmt")
	p.AddFunction(UTIL_FILE, "", s.SerializerMethod(), s.unionSerializer())
	p.AddImport(UTIL_FILE, "io")
	p.AddImport(UTIL_FILE, gogenavroImport)

	for _, f := range s.itemType {
		f.AddSerializer(p)
	}
}

func (s *unionField) AddDeserializer(p *generator.Package) {
	p.AddImport(UTIL_FILE, "fmt")
	p.AddFunction(UTIL_FILE, "", s.DeserializerMethod(), s.unionDeserializer())
	p.AddImport(UTIL_FILE, "io")
	p.AddImport(UTIL_FILE, gogenavroImport)

	for _, f := range s.itemType {
		f.AddDeserializer(p)
	}
}

func (s *unionField) ResolveReferences(n *Namespace) error {
	var err error
	for _, f := range s.itemType {
		err = f.ResolveReferences(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *unionField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	for i, item := range s.itemType {
		s.definition[i], err = item.Definition(scope)
		if err != nil {
			return nil, err
		}
	}
	return s.definition, nil
}

func (s *unionField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	lvalue = fmt.Sprintf("%v.%v", lvalue, s.itemType[0].Name())
	return s.itemType[0].DefaultValue(lvalue, rvalue)
}

func (s *unionField) Skip(r io.Reader) error {
	field, err := ReadLong(r)
	if err != nil {
		return err
	}

	if field < 0 || int(field) >= len(s.itemType) {
		return fmt.Errorf("Invalid value for %v", field)
	}

	return s.itemType[field].Skip(r)
}
