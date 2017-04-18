package types

import (
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/alanctgardner/gogen-avro/generator"
)

const recordStructDefTemplate = `type %v struct {
%v
}
`

const recordSchemaTemplate = `func (r %v) Schema() string {
 return %v
}
`

const recordConstructorTemplate = `
	func %v %v {
		v := &%v{
			%v
		}
		%v
		return v
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
	name QualifiedName
	aliases []QualifiedName
	fields   []*Field
	metadata map[string]interface{}
}

func NewRecordDefinition(name QualifiedName, aliases []QualifiedName, fields []*Field, metadata map[string]interface{}) *RecordDefinition {
	return &RecordDefinition {
		name: name,
		aliases: aliases,
		fields: fields,
		metadata: metadata,
	}
}

func (r *RecordDefinition) AvroName() QualifiedName {
	return r.name
}

func (r *RecordDefinition) Name() string {
	return generator.ToPublicName(r.name.Name)
}

func (r *RecordDefinition) GoType() string {
	return fmt.Sprintf("*%v", r.Name())
}

func (r *RecordDefinition) Aliases() []QualifiedName {
	return r.aliases
}

func (r *RecordDefinition) structFields() string {
	var fieldDefinitions string
	for _, f := range r.fields {
		fieldDefinitions += fmt.Sprintf("%v %v\n", f.GoName(), f.Type().GoType())
	}
	return fieldDefinitions
}

func (r *RecordDefinition) fieldSerializers() string {
	serializerMethods := "var err error\n"
	for _, f := range r.fields {
		serializerMethods += fmt.Sprintf("err = %v(r.%v, w)\nif err != nil {return err}\n", f.Type().SerializerMethod(), f.GoName())
	}
	return serializerMethods
}

func (r *RecordDefinition) fieldDeserializers() string {
	deserializerMethods := ""
	for _, f := range r.fields {
		deserializerMethods += fmt.Sprintf("str.%v, err = %v(r)\nif err != nil {return nil, err}\n", f.GoName(), f.Type().DeserializerMethod())
	}
	return deserializerMethods
}

func (r *RecordDefinition) structDefinition() string {
	return fmt.Sprintf(recordStructDefTemplate, r.Name(), r.structFields())
}

func (r *RecordDefinition) serializerMethodDef() string {
	return fmt.Sprintf("func %v(r %v, w io.Writer) error {\n%v\nreturn nil\n}", r.SerializerMethod(), r.GoType(), r.fieldSerializers())
}

func (r *RecordDefinition) deserializerMethodDef() string {
	return fmt.Sprintf(recordStructDeserializerTemplate, r.DeserializerMethod(), r.GoType(), r.Name(), r.fieldDeserializers())
}

func (r *RecordDefinition) SerializerMethod() string {
	return fmt.Sprintf("write%v", r.Name())
}

func (r *RecordDefinition) DeserializerMethod() string {
	return fmt.Sprintf("read%v", r.Name())
}

func (r *RecordDefinition) publicDeserializerMethod() string {
	return fmt.Sprintf("Deserialize%v", r.Name())
}

func (r *RecordDefinition) publicSerializerMethodDef() string {
	return fmt.Sprintf(recordStructPublicSerializerTemplate, r.GoType(), r.SerializerMethod())
}

func (r *RecordDefinition) publicDeserializerMethodDef() string {
	return fmt.Sprintf(recordStructPublicDeserializerTemplate, r.publicDeserializerMethod(), r.GoType(), r.DeserializerMethod())
}

func (r *RecordDefinition) filename() string {
	return generator.ToSnake(r.Name()) + ".go"
}

func (r *RecordDefinition) schemaMethodDef() string {
	def := r.Definition(make(map[QualifiedName]interface{}))
	schemaJson, _ := json.Marshal(def)
	return fmt.Sprintf(recordSchemaTemplate, r.GoType(), strconv.Quote(string(schemaJson)))
}

func (r *RecordDefinition) AddStruct(p *generator.Package) {
	// Import guard, to avoid circular dependencies
	if !p.HasStruct(r.filename(), r.GoType()) {
		p.AddStruct(r.filename(), r.GoType(), r.structDefinition())
		p.AddFunction(r.filename(), r.GoType(), "Schema", r.schemaMethodDef())
		p.AddFunction(r.filename(), r.GoType(), r.ConstructorMethod(), r.ConstructorMethodDef())
		for _, f := range r.fields {
			f.Type().AddStruct(p)
		}
	}
}

func (r *RecordDefinition) AddSerializer(p *generator.Package) {
	// Import guard, to avoid circular dependencies
	if !p.HasFunction(UTIL_FILE, "", r.SerializerMethod()) {
		p.AddImport(r.filename(), "io")
		p.AddFunction(UTIL_FILE, "", r.SerializerMethod(), r.serializerMethodDef())
		p.AddFunction(r.filename(), r.GoType(), "Serialize", r.publicSerializerMethodDef())
		for _, f := range r.fields {
			f.Type().AddSerializer(p)
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
			f.Type().AddDeserializer(p)
		}
	}
}

func (r *RecordDefinition) ResolveReferences(n *Namespace) error {
	var err error
	for _, f := range r.fields {
		err = f.Type().ResolveReferences(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RecordDefinition) Definition(scope map[QualifiedName]interface{}) interface{} {
	if _, ok := scope[r.name]; ok {
		return r.name.String()
	}
	scope[r.name] = 1
	fields := make([]map[string]interface{}, 0)
	for _, f := range r.fields {
		fields = append(fields, f.Definition(scope))
	}

	r.metadata["fields"] = fields
	return r.metadata
}

func (r *RecordDefinition) ConstructorMethod() string {
	return fmt.Sprintf("New%v()", r.Name())
}

func (r *RecordDefinition) fieldConstructors() string {
	constructors := ""
	for _, f := range r.fields {
		if constructor, ok := getConstructableForType(f.Type()); ok {
			constructors += fmt.Sprintf("%v: %v,\n", f.GoName(), constructor.ConstructorMethod())
		}
	}
	return constructors
}

func (r *RecordDefinition) defaultValues() string {
	defaults := ""
	for _, f := range r.fields {
		if f.hasDef {
			defaults += f.Type().DefaultValue(fmt.Sprintf("v.%v", f.GoName()), f.Default()) + "\n"
		}
	}
	return defaults
}

func (r *RecordDefinition) ConstructorMethodDef() string {
	return fmt.Sprintf(recordConstructorTemplate, r.ConstructorMethod(), r.GoType(), r.Name(), r.fieldConstructors(), r.defaultValues())
}

func (r *RecordDefinition) FieldByName(name string) *Field {
	for _, f := range r.fields {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

func (r *RecordDefinition) DefaultValue(lvalue string, rvalue interface{}) string {
	items := rvalue.(map[string]interface{})
	fieldSetters := ""
	for k, v := range items {
		field := r.FieldByName(k)
		fieldSetters += field.Type().DefaultValue(fmt.Sprintf("%v.%v", lvalue, field.GoName()), v) + "\n"
	}
	return fieldSetters
}
