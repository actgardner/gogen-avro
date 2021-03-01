// gogen-avro's internal representation of Avro schemas
package schema

import (
	"fmt"
	"strings"

	"github.com/actgardner/gogen-avro/v7/generator"
)

// invalidFieldNames is a list of field names that conflict with hard-coded method names on
// generated structs. These are converted to `Field_<name>` in the resulting struct to avoid errors.
var invalidFieldNames = map[string]interface{}{
	"Schema":               true,
	"Serialize":            true,
	"SchemaName":           true,
	"MarshalJSON":          true,
	"UnmarshalJSON":        true,
	"AvroCRC64Fingerprint": true,
	"SetBoolean":           true,
	"SetInt":               true,
	"SetLong":              true,
	"SetFloat":             true,
	"SetDouble":            true,
	"SetBytes":             true,
	"SetString":            true,
	"Get":                  true,
	"SetDefault":           true,
	"AppendMap":            true,
	"AppendArray":          true,
	"NullField":            true,
	"Finalize":             true,
}

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

func (f *Field) Index() int {
	return f.index
}

func (f *Field) Doc() string {
	return strings.ReplaceAll(f.doc, "\n", " ")
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
	name := generator.ToPublicName(f.avroName)
	if _, ok := invalidFieldNames[name]; ok {
		return "Field_" + name
	}
	return name
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
