package generator

import (
	"fmt"
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

type unionField struct {
	name     string
	hasDef   bool
	itemType []field
}

func (s *unionField) Name() string {
	return toPublicName(s.name)
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

func (s *unionField) AuxStructs(aux map[string]string, imports map[string]string) {
	if _, ok := aux[s.unionEnumType()]; ok {
		return
	}
	imports["fmt"] = "import \"fmt\""
	aux[s.unionEnumType()] = s.unionEnumDef()
	aux[s.SerializerMethod()] = s.unionSerializer()
	aux[s.FieldType()] = s.unionTypeDef()
	for _, f := range s.itemType {
		f.AuxStructs(aux, imports)
	}
}

func (s *unionField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.FieldType())
}
