package schema

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/actgardner/gogen-avro/generator"
	"github.com/actgardner/gogen-avro/schema/templates"
)

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

func (s *UnionField) UnionEnumType() string {
	return fmt.Sprintf("%vTypeEnum", s.Name())
}

func (s *UnionField) ItemName(item AvroType) string {
	return s.UnionEnumType() + item.Name()
}

func (s *UnionField) ItemTypes() []AvroType {
	return s.itemType
}

func (s *UnionField) filename() string {
	return generator.ToSnake(s.Name()) + ".go"
}

func (s *UnionField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *UnionField) ItemConstructor(f AvroType) string {
	if constructor, ok := getConstructableForType(f); ok {
		return constructor.ConstructorMethod()
	}
	return ""
}

func (s *UnionField) structDefinition() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("union").Parse(templates.UnionTemplate)
	if err != nil {
		return "", err
	}
	err = t.Execute(buf, s)
	return buf.String(), err
}

func (s *UnionField) AddStruct(p *generator.Package, containers bool) error {
	def, err := s.structDefinition()
	if err != nil {
		return err
	}

	p.AddFile(s.filename(), def)

	for _, f := range s.itemType {
		if err := f.AddStruct(p, containers); err != nil {
			return err
		}
	}
	return nil
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
	defaultType := s.itemType[0]
	init := fmt.Sprintf("%v = %v\n", lvalue, s.ConstructorMethod())
	lvalue = fmt.Sprintf("%v.%v", lvalue, defaultType.Name())
	constructorCall := ""
	if constructor, ok := getConstructableForType(defaultType); ok {
		constructorCall = fmt.Sprintf("%v = %v\n", lvalue, constructor.ConstructorMethod())
	}
	assignment, err := defaultType.DefaultValue(lvalue, rvalue)
	return init + constructorCall + assignment, err
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

func (s *UnionField) SimpleName() string {
	return s.GoType()
}
