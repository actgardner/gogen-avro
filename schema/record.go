package schema

import (
	"encoding/json"
)

type RecordDefinition struct {
	generatorMetadata

	name     QualifiedName
	aliases  []QualifiedName
	fields   []*Field
	doc      string
	metadata map[string]interface{}
}

func NewRecordDefinition(name QualifiedName, aliases []QualifiedName, fields []*Field, doc string, metadata map[string]interface{}) *RecordDefinition {
	return &RecordDefinition{
		name:     name,
		aliases:  aliases,
		fields:   fields,
		doc:      doc,
		metadata: metadata,
	}
}

func (r *RecordDefinition) AvroName() QualifiedName {
	return r.name
}

func (r *RecordDefinition) Aliases() []QualifiedName {
	return r.aliases
}

func (r *RecordDefinition) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[r.name]; ok {
		return r.name.String(), nil
	}
	scope[r.name] = 1
	fields := make([]map[string]interface{}, 0)
	for _, f := range r.fields {
		def, err := f.Definition(scope)
		if err != nil {
			return nil, err
		}
		fields = append(fields, def)
	}

	r.metadata["fields"] = fields
	return r.metadata, nil
}

func (r *RecordDefinition) GetReaderField(writerField *Field) *Field {
	for _, f := range r.fields {
		if f.IsSameField(writerField) {
			return f
		}
	}
	return nil
}

func (r *RecordDefinition) FieldByName(field string) *Field {
	for _, f := range r.fields {
		if f.NameMatchesAliases(field) {
			return f
		}
	}
	return nil
}

func (r *RecordDefinition) Fields() []*Field {
	return r.fields
}

func (s *RecordDefinition) IsReadableBy(d Definition) bool {
	reader, ok := d.(*RecordDefinition)
	return ok && reader.name == s.name
}

func (s *RecordDefinition) Doc() string {
	return s.doc
}

func (s *RecordDefinition) Schema() (string, error) {
	def, err := s.Definition(make(map[QualifiedName]interface{}))
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(def)
	return string(jsonBytes), err
}

func (s *RecordDefinition) Children() []AvroType {
	children := make([]AvroType, len(s.fields))
	for i, field := range s.fields {
		children[i] = field.Type()
	}
	return children
}
