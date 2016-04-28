package generator

import (
	"fmt"
)

const recordStructDefTemplate = `type %v struct {
%v
}
`

const recordStructPublicSerializerTemplate = `
func (r %v) Serialize(w io.Writer) error {
	return %v(&r, w)
}
`

const recordStructDeserializerTemplate = `
func %v(r io.Reader) (*%v, error) {
	var str %v
	var err error
	%v
	return &str, nil
}
`

const recordStructPublicDeserializerTemplate = `
func %v(r io.Reader) (*%v, error) {
	return %v(r)
}
`

type recordDefinition struct {
	name   string
	fields []field
}

func (r *recordDefinition) GoType() string {
	return toPublicName(r.name)
}

func (r *recordDefinition) structFields() string {
	var fieldDefinitions string
	for _, f := range r.fields {
		fieldDefinitions += fmt.Sprintf("%v %v\n", f.Name(), f.GoType())
	}
	return fieldDefinitions
}

func (r *recordDefinition) fieldSerializers() string {
	serializerMethods := "var err error\n"
	for _, f := range r.fields {
		serializerMethods += fmt.Sprintf("err = %v(r.%v, w)\nif err != nil {return err}\n", f.SerializerMethod(), f.Name())
	}
	return serializerMethods
}

func (r *recordDefinition) fieldDeserializers() string {
	deserializerMethods := ""
	for _, f := range r.fields {
		deserializerMethods += fmt.Sprintf("str.%v, err = %v(r)\nif err != nil {return nil, err}\n", f.Name(), f.DeserializerMethod())
	}
	return deserializerMethods
}

func (r *recordDefinition) structDefinition() string {
	return fmt.Sprintf(recordStructDefTemplate, r.GoType(), r.structFields())
}

func (r *recordDefinition) serializerMethodDef() string {
	return fmt.Sprintf("func %v(r *%v, w io.Writer) error {\n%v\nreturn nil\n}", r.serializerMethod(), r.GoType(), r.fieldSerializers())
}

func (r *recordDefinition) deserializerMethodDef() string {
	return fmt.Sprintf(recordStructDeserializerTemplate, r.deserializerMethod(), r.GoType(), r.GoType(), r.fieldDeserializers())
}

func (r *recordDefinition) serializerMethod() string {
	return fmt.Sprintf("write%v", r.GoType())
}

func (r *recordDefinition) deserializerMethod() string {
	return fmt.Sprintf("read%v", r.GoType())
}

func (r *recordDefinition) publicDeserializerMethod() string {
	return fmt.Sprintf("Deserialize%v", r.GoType())
}

func (r *recordDefinition) publicSerializerMethodDef() string {
	return fmt.Sprintf(recordStructPublicSerializerTemplate, r.GoType(), r.serializerMethod())
}

func (r *recordDefinition) publicDeserializerMethodDef() string {
	return fmt.Sprintf(recordStructPublicDeserializerTemplate, r.publicDeserializerMethod(), r.GoType(), r.deserializerMethod())
}

func (r *recordDefinition) filename() string {
	return toSnake(r.GoType()) + ".go"
}

func (r *recordDefinition) AddStruct(p *Package) {
	// Import guard, to avoid circular dependencies
	if !p.hasStruct(r.filename(), r.GoType()) {
		p.addStruct(r.filename(), r.GoType(), r.structDefinition())
		for _, f := range r.fields {
			f.AddStruct(p)
		}
	}
}

func (r *recordDefinition) AddSerializer(p *Package) {
	// Import guard, to avoid circular dependencies
	if !p.hasFunction(r.filename(), "", r.serializerMethod()) {
		p.addImport(r.filename(), "io")
		p.addFunction(UTIL_FILE, "", r.serializerMethod(), r.serializerMethodDef())
		p.addFunction(r.filename(), r.GoType(), "Serialize", r.publicSerializerMethodDef())
		for _, f := range r.fields {
			f.AddSerializer(p)
		}
	}
}

func (r *recordDefinition) AddDeserializer(p *Package) {
	// Import guard, to avoid circular dependencies
	if !p.hasFunction(r.filename(), "", r.deserializerMethod()) {
		p.addImport(r.filename(), "io")
		p.addFunction(UTIL_FILE, "", r.deserializerMethod(), r.deserializerMethodDef())
		p.addFunction(r.filename(), "", r.publicDeserializerMethod(), r.publicDeserializerMethodDef())
		for _, f := range r.fields {
			f.AddDeserializer(p)
		}
	}
}
