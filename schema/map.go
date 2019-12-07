package schema

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/actgardner/gogen-avro/generator"
	"github.com/actgardner/gogen-avro/schema/templates"
)

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

func (s *MapField) filename() string {
	return generator.ToSnake(s.Name()) + ".go"
}

func (s *MapField) AddStruct(p *generator.Package, containers bool) error {
	def, err := s.structDefinition()
	if err != nil {
		return err
	}

	p.AddFile(s.filename(), def)

	return s.itemType.AddStruct(p, containers)
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

func (s *MapField) structDefinition() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("map").Parse(templates.MapTemplate)
	if err != nil {
		return "", err
	}
	err = t.Execute(buf, s)
	return buf.String(), err
}

func (s *MapField) SimpleName() string {
	return s.Name()
}

func (s *MapField) ItemConstructable() string {
	if constructor, ok := getConstructableForType(s.itemType); ok {
		return fmt.Sprintf("v = %v\n", constructor.ConstructorMethod())
	}
	return ""
}
