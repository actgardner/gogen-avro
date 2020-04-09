// gogen-avro's internal representation of Avro schemas
package schema

import (
	"fmt"

	"github.com/actgardner/gogen-avro/generator"
)

type Field struct {
	avroName   string
	avroType   AvroType
	defValue   interface{}
	aliases    []string
	hasDef     bool
	doc        string
	definition map[string]interface{}
	fieldTags  string
	index      int
}

func NewField(avroName string, avroType AvroType, defValue interface{}, hasDef bool, aliases []string, doc string, definition map[string]interface{}, index int, fieldTags string) *Field {
	return &Field{
		avroName:   avroName,
		avroType:   avroType,
		defValue:   defValue,
		hasDef:     hasDef,
		aliases:    aliases,
		doc:        doc,
		definition: definition,
		fieldTags:  fieldTags,
		index:      index,
	}
}

func (f *Field) Name() string {
	return f.avroName
}

func (f *Field) SimpleName() string {
	return generator.ToPublicSimpleName(f.avroName)
}

func (f *Field) Index() int {
	return f.index
}

func (f *Field) Doc() string {
	return f.doc
}

// Tags returns a field go struct tags if defined.
func (f *Field) Tags() string {
	jsonTag := fmt.Sprintf("json:%q", f.avroName)
	if f.fieldTags == "" {
		return jsonTag
	}
	return f.fieldTags + " " + jsonTag
}

func (f *Field) GoName() string {
	return generator.ToPublicName(f.avroName)
}

// IsSameField checks whether two fields have the same name or any of their aliases are the same, in which case they're the same for purposes of schema evolution
func (f *Field) IsSameField(otherField *Field) bool {
	if otherField.NameMatchesAliases(f.avroName) {
		return true
	}

	for _, n := range f.aliases {
		if otherField.NameMatchesAliases(n) {
			return true
		}
	}

	return false
}

func (f *Field) NameMatchesAliases(name string) bool {
	if name == f.avroName {
		return true
	}

	for _, n := range f.aliases {
		if n == name {
			return true
		}
	}

	return false
}

func (f *Field) HasDefault() bool {
	return f.hasDef
}

func (f *Field) Default() interface{} {
	return f.defValue
}

func (f *Field) Type() AvroType {
	if f == nil {
		return nil
	}
	return f.avroType
}

func (f *Field) Definition(scope map[QualifiedName]interface{}) (map[string]interface{}, error) {
	def := copyDefinition(f.definition)
	var err error
	def["type"], err = f.avroType.Definition(scope)
	if err != nil {
		return nil, err
	}
	return def, nil
}
