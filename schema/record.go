package schema

import (
	"encoding/json"
	"fmt"

	"github.com/actgardner/gogen-avro/generator"
)

type RecordDefinition struct {
	qualifiedField
	fields   []*Field
	doc      string
	children []AvroType
}

func NewRecordDefinition(qname QualifiedName, aliases []QualifiedName, fields []*Field, doc string, definition interface{}) *RecordDefinition {
	r := &RecordDefinition{
		qualifiedField: newQualifiedField(qname, aliases, definition),
		fields:         fields,
		doc:            doc,
	}
	r.children = makeChildren(fields)
	// Override goType and simpleName due to namespaces
	r.setSimpleName(r.GoType())
	r.setGoType(generator.ToPublicName(qname.String()))
	return r
}

// Convenience function for caching children
func makeChildren(fields []*Field) []AvroType {
	children := make([]AvroType, len(fields))
	for i, field := range fields {
		children[i] = field.Type()
	}
	return children
}

func (r *RecordDefinition) Fields() []*Field {
	return r.fields
}

func (s *RecordDefinition) Doc() string {
	return s.doc
}

func (s *RecordDefinition) Children() []AvroType {
	return s.children
}

func (s *RecordDefinition) Schema() (string, error) {
	def0, err := s.Definition(make(map[QualifiedName]interface{}))
	if err != nil {
		return "", err
	}
	def := def0.(map[string]interface{})
	delete(def, "namespace")
	def["name"] = s.qname.String()
	jsonBytes, err := json.Marshal(def)
	return string(jsonBytes), err
}

func (r *RecordDefinition) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[r.qname]; ok {
		return r.qname.String(), nil
	}
	scope[r.qname] = 1
	enhancedDef := copyDefinition(r.definitionAsStringMap())
	fields := make([]map[string]interface{}, 0)
	for _, f := range r.fields {
		def, err := f.Definition(scope)
		if err != nil {
			return nil, err
		}
		fields = append(fields, def)
	}

	enhancedDef["fields"] = fields
	return enhancedDef, nil
}

func (r *RecordDefinition) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	items := rvalue.(map[string]interface{})
	fieldSetters := ""
	for k, v := range items {
		field := r.fieldByName(k)
		fieldSetter, err := field.Type().DefaultValue(fmt.Sprintf("%v.%v", lvalue, field.GoType()), v)
		if err != nil {
			return "", err
		}

		fieldSetters += fieldSetter + "\n"
	}
	return fieldSetters, nil
}

func (s *RecordDefinition) IsReadableBy(d AvroType, visited map[QualifiedName]interface{}) bool {
	// If there's a circular reference, don't evaluate every field on the second pass
	if _, ok := visited[s.qname]; ok {
		return true
	}

	// Reader not a Record, maybe a union?
	if _, ok := d.(*RecordDefinition); !ok {
		if u, ok := d.(*UnionField); ok {
			return isReadableByUnion(s, u, visited)
		}
		return false
	}

	visited[s.qname] = true
	reader := d.(*RecordDefinition)
	for _, readerField := range reader.Fields() {
		writerField := s.GetReaderField(readerField)
		// Two schemas are incompatible if the reader has a field with no default value that is not present in the writer schema
		if writerField == nil && !readerField.HasDefault() {
			return false
		}

		// The two schemas are incompatible if two fields with the same name have different schemas
		if writerField != nil && !writerField.Type().IsReadableBy(readerField.Type(), visited) {
			return false
		}

	}
	return true
}

func (r *RecordDefinition) ConstructorMethod() string {
	return fmt.Sprintf("%s{}", r.Name())
}

func (r *RecordDefinition) ConstructableForField(f *Field) string {
	if constructor, ok := getConstructableForType(f.Type()); ok {
		return fmt.Sprintf("r.%v = %v\n", f.GoType(), constructor.ConstructorMethod())
	}
	return ""
}

func (r *RecordDefinition) DefaultForField(f *Field) (string, error) {
	if f.Type().IsOptional() && f.Default() == nil {
		return fmt.Sprintf("r.%v = nil", f.GoType()), nil
	}
	return f.Type().DefaultValue(fmt.Sprintf("r.%v", f.GoType()), f.Default())
}

func (r *RecordDefinition) GetReaderField(writerField *Field) *Field {
	for _, f := range r.fields {
		if f.isSameField(writerField) {
			return f
		}
	}
	return nil
}

func (r *RecordDefinition) fieldByName(field string) *Field {
	for _, f := range r.fields {
		if f.nameMatchesAliases(field) {
			return f
		}
	}
	return nil
}
