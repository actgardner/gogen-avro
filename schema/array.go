package schema

import (
	"fmt"
)

type ArrayField struct {
	qualifiedField
	itemType AvroType
}

func NewArrayField(itemType AvroType, definition map[string]interface{}) *ArrayField {
	a := &ArrayField{}
	a.definition = definition
	a.setItemType(itemType)

	// If the item type is a reference, its final gotype will not be known until resolution
	if ref, ok := itemType.(*Reference); ok {
		ref.AddResolver(a)
	}
	return a
}

func createArrayQName(name string) QualifiedName {
	return QualifiedName{Name: fmt.Sprintf("Array%s", name)}
}

// Resolve runtime data when item type is a reference, since its data
// cannot be always known at this array's creation time.
func (s *ArrayField) Resolve(ref Reference) {
	s.setItemType(ref.refType)
}

func (s *ArrayField) ItemType() AvroType {
	return s.itemType
}

func (s *ArrayField) setItemType(itemType AvroType) {
	s.itemType = itemType
	s.setQualifiedName(QualifiedName{Name: "Array" + itemType.Name()})
	s.setGoType(fmt.Sprintf("[]%s", itemType.GoType()))
}

func (s *ArrayField) Children() []AvroType {
	return []AvroType{s.itemType}
}

func (s *ArrayField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	def := copyDefinition(s.definitionAsStringMap())
	var err error
	def["items"], err = s.itemType.Definition(scope)
	return def, err
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
	return s.Name()
}

func (s *ArrayField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	if reader, ok := f.(*ArrayField); ok {
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

func (s *ArrayField) ItemConstructable() string {
	if constructor, ok := getConstructableForType(s.itemType); ok {
		return fmt.Sprintf("v = %v\n", constructor.ConstructorMethod())
	}
	return ""
}
