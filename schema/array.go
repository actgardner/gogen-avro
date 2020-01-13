package schema

import (
	"fmt"

	"github.com/actgardner/gogen-avro/v7/generator"
)

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

func (r *ArrayField) filename() string {
	return generator.ToSnake(r.Name()) + ".go"
}

func (s *ArrayField) GoType() string {
	return fmt.Sprintf("[]%v", s.itemType.GoType())
}

func (s *ArrayField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Name())
}

func (s *ArrayField) ItemType() AvroType {
	return s.itemType
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

func (s *ArrayField) ItemConstructable() string {
	if constructor, ok := getConstructableForType(s.itemType); ok {
		return fmt.Sprintf("v = %v\n", constructor.ConstructorMethod())
	}
	return ""
}

func (s *ArrayField) Children() []AvroType {
	return []AvroType{s.itemType}
}
