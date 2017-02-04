package types

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
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
	return fmt.Errorf("Invalid value for %v")
}
`

const unionDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	field, err := readLong(r)
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
	name         string
	hasDefault   bool
	defaultValue interface{}
	itemType     []Field
}

func (s *unionField) HasDefault() bool {
	return s.hasDefault
}

func (s *unionField) Default() interface{} {
	return s.defaultValue
}

func (s *unionField) AvroName() string {
	return s.name
}

func (s *unionField) GoName() string {
	return generator.ToPublicName(s.name)
}

func (s *unionField) FieldType() string {
	var unionFields = "Union"
	for _, i := range s.itemType {
		unionFields += i.FieldType()
	}
	return unionFields
}

func (s *unionField) GoType() string {
	return s.FieldType()
}

func (s *unionField) unionEnumType() string {
	return fmt.Sprintf("%vTypeEnum", s.FieldType())
}

func (s *unionField) unionEnumDef() string {
	var unionTypes string
	for i, item := range s.itemType {
		unionTypes += fmt.Sprintf("%v %v = %v\n", s.unionEnumType()+item.FieldType(), s.unionEnumType(), i)
	}
	return fmt.Sprintf("type %v int\nconst(\n%v)\n", s.unionEnumType(), unionTypes)
}

func (s *unionField) unionTypeDef() string {
	var unionFields string
	for _, i := range s.itemType {
		unionFields += fmt.Sprintf("%v %v\n", i.FieldType(), i.GoType())
	}
	unionFields += fmt.Sprintf("UnionType %v", s.unionEnumType())
	return fmt.Sprintf("type %v struct{\n%v\n}\n", s.FieldType(), unionFields)
}

func (s *unionField) unionSerializer() string {
	switchCase := ""
	for _, t := range s.itemType {
		switchCase += fmt.Sprintf("case %v:\nreturn %v(r.%v, w)\n", s.unionEnumType()+t.FieldType(), t.SerializerMethod(), t.FieldType())
	}
	return fmt.Sprintf(unionSerializerTemplate, s.SerializerMethod(), s.GoType(), switchCase, s.GoType())
}

func (s *unionField) unionDeserializer() string {
	switchCase := ""
	for _, t := range s.itemType {
		switchCase += fmt.Sprintf("case %v:\nval, err :=  %v(r)\nif err != nil {return unionStr, err}\nunionStr.%v = val\n", s.unionEnumType()+t.FieldType(), t.DeserializerMethod(), t.FieldType())
	}
	return fmt.Sprintf(unionDeserializerTemplate, s.DeserializerMethod(), s.GoType(), s.GoType(), s.unionEnumType(), switchCase, s.GoType())
}

func (s *unionField) filename() string {
	return generator.ToSnake(s.GoType()) + ".go"
}

func (s *unionField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.FieldType())
}

func (s *unionField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.FieldType())
}

func (s *unionField) AddStruct(p *generator.Package) {
	p.AddStruct(s.filename(), s.unionEnumType(), s.unionEnumDef())
	p.AddStruct(s.filename(), s.FieldType(), s.unionTypeDef())
	for _, f := range s.itemType {
		f.AddStruct(p)
	}
}

func (s *unionField) AddSerializer(p *generator.Package) {
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

func (s *unionField) AddDeserializer(p *generator.Package) {
	p.AddImport(UTIL_FILE, "fmt")
	p.AddFunction(UTIL_FILE, "", s.DeserializerMethod(), s.unionDeserializer())
	p.AddFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.AddImport(UTIL_FILE, "io")
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

func (s *unionField) Schema(names map[QualifiedName]interface{}) interface{} {
	unionDefs := make([]interface{}, 0, len(s.itemType))
	for _, item := range s.itemType {
		unionDefs = append(unionDefs, item.Schema(names))
	}
	return unionDefs
}
