// gogen-avro's internal representation of Avro schemas
package schema

import (
	"github.com/actgardner/gogen-avro/generator"
)

type Field struct {
	abstractField
	aliases   []string
	doc       string
	fieldType AvroType
	fieldTags string
	index     int
}

func NewField(name string, fieldType AvroType, aliases []string, doc string, definition map[string]interface{}, index int, fieldTags string) *Field {
	return &Field{
		abstractField: abstractField{
			name:       name,
			simpleName: generator.ToPublicSimpleName(name),
			goType:     generator.ToPublicName(name),
			definition: definition,
		},
		aliases:   aliases,
		doc:       doc,
		fieldType: fieldType,
		fieldTags: fieldTags,
		index:     index,
	}
}

func (f *Field) Index() int {
	return f.index
}

func (f *Field) Doc() string {
	return f.doc
}

// Tags returns a field go struct tags if defined.
func (f *Field) Tags() string {
	return f.fieldTags
}

func (f *Field) HasDefault() bool {
	_, ok := f.definitionAsStringMap()["default"]
	return ok
}

func (f *Field) Default() interface{} {
	return f.definitionAsStringMap()["default"]
}

func (f *Field) Type() AvroType {
	if f == nil {
		return nil
	}
	return f.fieldType
}

func (f *Field) Definition(scope map[QualifiedName]interface{}) (map[string]interface{}, error) {
	def := copyDefinition(f.definitionAsStringMap())
	var err error
	def["type"], err = f.fieldType.Definition(scope)
	if err != nil {
		return nil, err
	}
	return def, nil
}

// isSameField checks whether two fields have the same name
// or any of their aliases are the same, in which case they're
// the same for purposes of schema evolution
func (f *Field) isSameField(otherField *Field) bool {
	if otherField.nameMatchesAliases(f.name) {
		return true
	}

	for _, n := range f.aliases {
		if otherField.nameMatchesAliases(n) {
			return true
		}
	}

	return false
}

func (f *Field) nameMatchesAliases(name string) bool {
	if name == f.name {
		return true
	}

	for _, n := range f.aliases {
		if n == name {
			return true
		}
	}

	return false
}
