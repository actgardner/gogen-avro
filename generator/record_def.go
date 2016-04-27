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

func (r *recordDefinition) structDefinition() string {
	return fmt.Sprintf(recordStructDefTemplate, r.GoType(), r.structFields())
}

func (r *recordDefinition) serializerMethodDef() string {
	return fmt.Sprintf("func %v(r *%v, w io.Writer) error {\n%v\nreturn nil\n}", r.serializerMethod(), r.GoType(), r.fieldSerializers())
}

func (r *recordDefinition) serializerMethod() string {
	return fmt.Sprintf("write%v", r.GoType())
}

func (r *recordDefinition) publicSerializerMethodDef() string {
	return fmt.Sprintf(recordStructPublicSerializerTemplate, r.GoType(), r.serializerMethod())
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
