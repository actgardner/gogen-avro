package schema

type RecordDefinition struct {
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

func (r *RecordDefinition) Children() []AvroType {
	children := make([]AvroType, len(r.fields))
	for i, f := range r.fields {
		children[i] = f.avroType
	}
	return children
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
