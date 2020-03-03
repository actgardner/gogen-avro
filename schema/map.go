package schema

import (
	"fmt"
)

type MapField struct {
	primitiveField
	itemType AvroType
}

func NewMapField(itemType AvroType, definition map[string]interface{}) *MapField {
	name := "Map" + itemType.Name()
	return &MapField{
		primitiveField: newPrimitiveField(name, name, definition, "write"+name),
		itemType:       itemType,
	}
}

func (s *MapField) ItemType() AvroType {
	return s.itemType
}

func (s *MapField) Children() []AvroType {
	return []AvroType{s.itemType}
}

func (s *MapField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	def := copyDefinition(s.definitionAsStringMap())
	var err error
	def["values"], err = s.itemType.Definition(scope)
	if err != nil {
		return nil, err
	}
	return def, nil
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

func (s *MapField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	if reader, ok := f.(*MapField); ok {
		return s.ItemType().IsReadableBy(reader.ItemType(), visited)
	}
	if s.primitiveField.IsReadableBy(f, visited) {
		return true
	}
	if union, ok := f.(*UnionField); ok {
		return isReadableByUnion(s, union, visited)
	}
	return false
}

func (s *MapField) ItemConstructable() string {
	if constructor, ok := getConstructableForType(s.itemType); ok {
		return fmt.Sprintf("v = %v\n", constructor.ConstructorMethod())
	}
	return ""
}

func (s *MapField) ConstructorMethod() string {
	return fmt.Sprintf("New%s()", s.Name())
}
