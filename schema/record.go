package schema

import (
	"encoding/json"
	"fmt"
	"github.com/actgardner/gogen-avro/generator"
	"strconv"
)

const recordStructDefTemplate = `
%v
type %v struct {
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

const recordWriterTemplate = `
func %v(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := &%v{}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}
`

const recordFieldTemplate = `
func (_ %[1]v) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ %[1]v) SetInt(v int32) { panic("Unsupported operation") }
func (_ %[1]v) SetLong(v int64) { panic("Unsupported operation") }
func (_ %[1]v) SetFloat(v float32) { panic("Unsupported operation") }
func (_ %[1]v) SetDouble(v float64) { panic("Unsupported operation") }
func (_ %[1]v) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ %[1]v) SetString(v string) { panic("Unsupported operation") }
func (r %[1]v) Get(i int) types.Field {
	switch (i) {
		%[2]v
	}
	panic("Unknown field index")
}
func (_ %[1]v) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ %[1]v) AppendArray() types.Field { panic("Unsupported operation") }
`

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
		if f.Doc() != "" {
			fieldDefinitions += fmt.Sprintf("\n// %v\n", f.Doc())
		}
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

func (r *RecordDefinition) structDefinition() string {
	var doc string
	if r.doc != "" {
		doc = fmt.Sprintf("// %v", r.doc)
	}
	return fmt.Sprintf(recordStructDefTemplate, doc, r.Name(), r.structFields())
}

func (r *RecordDefinition) serializerMethodDef() string {
	return fmt.Sprintf("func %v(r %v, w io.Writer) error {\n%v\nreturn nil\n}", r.SerializerMethod(), r.GoType(), r.fieldSerializers())
}

func (r *RecordDefinition) SerializerMethod() string {
	return fmt.Sprintf("write%v", r.Name())
}

func (r *RecordDefinition) recordWriterMethod() string {
	return fmt.Sprintf("New%vWriter", r.Name())
}

func (r *RecordDefinition) recordWriterMethodDef() string {
	return fmt.Sprintf(recordWriterTemplate, r.recordWriterMethod(), r.Name())
}

func (r *RecordDefinition) publicSerializerMethodDef() string {
	return fmt.Sprintf(recordStructPublicSerializerTemplate, r.GoType(), r.SerializerMethod())
}

func (r *RecordDefinition) filename() string {
	return generator.ToSnake(r.Name()) + ".go"
}

func (r *RecordDefinition) schemaMethodDef() (string, error) {
	def, err := r.Definition(make(map[QualifiedName]interface{}))
	if err != nil {
		return "", err
	}

	schemaJson, _ := json.Marshal(def)
	return fmt.Sprintf(recordSchemaTemplate, r.GoType(), strconv.Quote(string(schemaJson))), nil
}

func (r *RecordDefinition) AddStruct(p *generator.Package, containers bool) error {
	// Import guard, to avoid circular dependencies
	if !p.HasStruct(r.filename(), r.GoType()) {
		p.AddStruct(r.filename(), r.GoType(), r.structDefinition())
		schemaDef, err := r.schemaMethodDef()
		if err != nil {
			return err
		}

		p.AddFunction(r.filename(), r.GoType(), "Schema", schemaDef)
		constructorMethodDef, err := r.ConstructorMethodDef()
		if err != nil {
			return err
		}

		if containers {
			p.AddImport(r.filename(), "github.com/actgardner/gogen-avro/container")
			p.AddFunction(r.filename(), "", r.recordWriterMethod(), r.recordWriterMethodDef())
		}

		p.AddImport(r.filename(), "github.com/actgardner/gogen-avro/types")
		p.AddFunction(r.filename(), r.GoType(), "fieldTemplate", r.FieldsMethodDef())
		p.AddFunction(r.filename(), r.GoType(), r.ConstructorMethod(), constructorMethodDef)
		for _, f := range r.fields {
			f.Type().AddStruct(p, containers)
		}
	}
	return nil
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

func (r *RecordDefinition) ConstructorMethod() string {
	return fmt.Sprintf("New%v()", r.Name())
}

func (r *RecordDefinition) fieldConstructors() (string, error) {
	constructors := ""
	for _, f := range r.fields {
		if constructor, ok := getConstructableForType(f.Type()); ok {
			constructors += fmt.Sprintf("%v: %v,\n", f.GoName(), constructor.ConstructorMethod())
		}
	}
	return constructors, nil
}

func (r *RecordDefinition) defaultValues() (string, error) {
	defaults := ""
	for _, f := range r.fields {
		if f.hasDef {
			def, err := f.Type().DefaultValue(fmt.Sprintf("v.%v", f.GoName()), f.Default())
			if err != nil {
				return "", err
			}
			defaults += def + "\n"
		}
	}
	return defaults, nil
}

func (r *RecordDefinition) FieldsMethodDef() string {
	getBody := ""
	for i, f := range r.fields {
		if f.Type().WrapperType() == "" {
			getBody += fmt.Sprintf("case %v:\nreturn r.%v\nbreak\n", i, f.Name())
		} else {
			getBody += fmt.Sprintf("case %v:\nreturn (*%v)(&r.%v)\nbreak\n", i, f.Type().WrapperType(), f.Name())
		}
	}
	return fmt.Sprintf(recordFieldTemplate, r.GoType(), getBody)
}

func (r *RecordDefinition) ConstructorMethodDef() (string, error) {
	defaults, err := r.defaultValues()
	if err != nil {
		return "", err
	}

	fieldConstructors, err := r.fieldConstructors()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(recordConstructorTemplate, r.ConstructorMethod(), r.GoType(), r.Name(), fieldConstructors, defaults), nil
}

func (r *RecordDefinition) FieldByName(name string) *Field {
	for _, f := range r.fields {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

func (r *RecordDefinition) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	items := rvalue.(map[string]interface{})
	fieldSetters := ""
	for k, v := range items {
		field := r.FieldByName(k)
		fieldSetter, err := field.Type().DefaultValue(fmt.Sprintf("%v.%v", lvalue, field.GoName()), v)
		if err != nil {
			return "", err
		}

		fieldSetters += fieldSetter + "\n"
	}
	return fieldSetters, nil
}

func (r *RecordDefinition) Fields() []*Field {
	return r.fields
}

func (s *RecordDefinition) IsReadableBy(d Definition) bool {
	reader, ok := d.(*RecordDefinition)
	if !ok {
		return false
	}

reader:
	for _, readerField := range reader.fields {
		for _, writerField := range s.fields {
			if readerField.Name() == writerField.Name() {
				if !writerField.Type().IsReadableBy(readerField.Type()) {
					return false
				}
				continue reader
			}
		}
		if !readerField.HasDefault() {
			return false
		}
	}
	return true
}
