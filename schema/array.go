package schema

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/actgardner/gogen-avro/generator"
	"github.com/actgardner/gogen-avro/schema/templates"
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

func (s *ArrayField) ItemType() AvroType {
	return s.itemType
}

func (s *ArrayField) ResolveReferences(n *Namespace) error {
	return s.itemType.ResolveReferences(n)
}

func (s *ArrayField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	s.definition["items"], err = s.itemType.Definition(scope)
	if err != nil {
		return nil, err
	}

	return s.definition, nil
}

func (s *ArrayField) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*ArrayField); ok {
		return s.ItemType().IsReadableBy(reader.ItemType())
	}
	return false
}
