package generator

import (
	"fmt"
)

const recordStructDefTemplate = `type %v struct {
%v
}

func (r %v) Serialize(w io.Writer) error {
	return write%v(r, w)
}
`

type recordDefinition struct {
	name   string
	fields []field
}

func (r *recordDefinition) goName() string {
	return toPublicName(r.name)
}

func (r *recordDefinition) structFields() string {
	var fieldDefinitions string
	for _, f := range r.fields {
		fieldDefinitions += fmt.Sprintf("%v %v\n", f.Name(), f.GoType())
	}
	return fieldDefinitions
}

/* Get the import and namespace maps for this record */
func (r *recordDefinition) namespaceMap(imports map[string]string, ns map[string]string) {
	imports["io"] = "import \"io\""
	ns[r.serializerMethod()] = r.serializerMethodDef()
	for _, f := range r.fields {
		f.SerializerNs(imports, ns)
	}
}

func (r *recordDefinition) fieldSerializers() string {
	serializerMethods := "var err error\n"
	for _, f := range r.fields {
		serializerMethods += fmt.Sprintf("err = %v(r.%v, w)\nif err != nil {return err}\n", f.SerializerMethod(), f.Name())
	}
	return serializerMethods
}

func (r *recordDefinition) structDefinition() string {
	return fmt.Sprintf(recordStructDefTemplate, r.goName(), r.structFields(), r.goName(), r.goName())
}

func (r *recordDefinition) serializerMethodDef() string {
	return fmt.Sprintf("func %v(r %v, w io.Writer) error {\n%v\nreturn nil\n}", r.serializerMethod(), r.goName(), r.fieldSerializers())
}

func (r *recordDefinition) serializerMethod() string {
	return fmt.Sprintf("write%v", r.goName())
}
