package types

import (
	"encoding/json"
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
	"strconv"
)

const recordStructDefTemplate = `type %v struct {
%v
}
`

const recordSchemaTemplate = `func (r %v) Schema() string {
 return %v
}
`

const recordStructPublicSerializerTemplate = `
func (r %v) Serialize(w io.Writer) error {
	return %v(r, w)
}
`

const recordStructDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	var str = &%v{}
	var err error
	%v
	return str, nil
}
`

const recordStructPublicDeserializerTemplate = `
func %v(r io.Reader) (%v, error) {
	return %v(r)
}
`

type RecordDefinition struct {
	name    QualifiedName
	aliases []QualifiedName
	fields  []Field
}

func (r *RecordDefinition) AvroName() QualifiedName {
	return r.name
}

func (r *RecordDefinition) Aliases() []QualifiedName {
	return r.aliases
}

func (r *RecordDefinition) GoType() string {
	return fmt.Sprintf("*%v", r.FieldType())
}

func (r *RecordDefinition) FieldType() string {
	return generator.ToPublicName(r.name.Name)
}

func (r *RecordDefinition) structFields() string {
	var fieldDefinitions string
	for _, f := range r.fields {
		fieldDefinitions += fmt.Sprintf("%v %v\n", f.Name(), f.GoType())
	}
	return fieldDefinitions
}

func (r *RecordDefinition) fieldSerializers() string {
	serializerMethods := "var err error\n"
	for _, f := range r.fields {
		serializerMethods += fmt.Sprintf("err = %v(r.%v, w)\nif err != nil {return err}\n", f.SerializerMethod(), f.Name())
	}
	return serializerMethods
}

func (r *RecordDefinition) fieldDeserializers() string {
	deserializerMethods := ""
	for _, f := range r.fields {
		deserializerMethods += fmt.Sprintf("str.%v, err = %v(r)\nif err != nil {return nil, err}\n", f.Name(), f.DeserializerMethod())
	}
	return deserializerMethods
}

func (r *RecordDefinition) structDefinition() string {
	return fmt.Sprintf(recordStructDefTemplate, r.FieldType(), r.structFields())
}

func (r *RecordDefinition) serializerMethodDef() string {
	return fmt.Sprintf("func %v(r %v, w io.Writer) error {\n%v\nreturn nil\n}", r.SerializerMethod(), r.GoType(), r.fieldSerializers())
}

func (r *RecordDefinition) deserializerMethodDef() string {
	return fmt.Sprintf(recordStructDeserializerTemplate, r.DeserializerMethod(), r.GoType(), r.FieldType(), r.fieldDeserializers())
}

func (r *RecordDefinition) SerializerMethod() string {
	return fmt.Sprintf("write%v", r.FieldType())
}

func (r *RecordDefinition) DeserializerMethod() string {
	return fmt.Sprintf("read%v", r.FieldType())
}

func (r *RecordDefinition) publicDeserializerMethod() string {
	return fmt.Sprintf("Deserialize%v", r.FieldType())
}

func (r *RecordDefinition) publicSerializerMethodDef() string {
	return fmt.Sprintf(recordStructPublicSerializerTemplate, r.GoType(), r.SerializerMethod())
}

func (r *RecordDefinition) publicDeserializerMethodDef() string {
	return fmt.Sprintf(recordStructPublicDeserializerTemplate, r.publicDeserializerMethod(), r.GoType(), r.DeserializerMethod())
}

func (r *RecordDefinition) filename() string {
	return generator.ToSnake(r.FieldType()) + ".go"
}

func (r *RecordDefinition) schemaMethod() string {
	schemaJson, _ := json.Marshal(r.Schema(make(map[QualifiedName]interface{})))
	return fmt.Sprintf(recordSchemaTemplate, r.GoType(), strconv.Quote(string(schemaJson)))
}

func (r *RecordDefinition) AddStruct(p *generator.Package) {
	// Import guard, to avoid circular dependencies
	if !p.HasStruct(r.filename(), r.GoType()) {
		p.AddStruct(r.filename(), r.GoType(), r.structDefinition())
		for _, f := range r.fields {
			f.AddStruct(p)
		}
		p.AddFunction(r.filename(), r.GoType(), "Schema", r.schemaMethod())
	}
}

func (r *RecordDefinition) AddSerializer(p *generator.Package) {
	// Import guard, to avoid circular dependencies
	if !p.HasFunction(UTIL_FILE, "", r.SerializerMethod()) {
		p.AddImport(r.filename(), "io")
		p.AddFunction(UTIL_FILE, "", r.SerializerMethod(), r.serializerMethodDef())
		p.AddFunction(r.filename(), r.GoType(), "Serialize", r.publicSerializerMethodDef())
		for _, f := range r.fields {
			f.AddSerializer(p)
		}
	}
}

func (r *RecordDefinition) AddDeserializer(p *generator.Package) {
	// Import guard, to avoid circular dependencies
	if !p.HasFunction(UTIL_FILE, "", r.DeserializerMethod()) {
		p.AddImport(r.filename(), "io")
		p.AddFunction(UTIL_FILE, "", r.DeserializerMethod(), r.deserializerMethodDef())
		p.AddFunction(r.filename(), "", r.publicDeserializerMethod(), r.publicDeserializerMethodDef())
		for _, f := range r.fields {
			f.AddDeserializer(p)
		}
	}
}

func (r *RecordDefinition) ResolveReferences(n *Namespace) error {
	var err error
	for _, f := range r.fields {
		err = f.ResolveReferences(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RecordDefinition) Schema(names map[QualifiedName]interface{}) interface{} {
	name := r.name.String()
	if _, ok := names[r.name]; ok {
		return name
	}
	names[r.name] = 1
	fields := make([]interface{}, 0, len(r.fields))
	for _, f := range r.fields {
		fieldDef := map[string]interface{}{
			"name": f.Name(),
			"type": f.Schema(names),
		}
		if f.HasDefault() {
			fieldDef["default"] = f.Default()
		}
		fields = append(fields, fieldDef)
	}
	return map[string]interface{}{
		"type":   "record",
		"name":   name,
		"fields": fields,
	}
}
