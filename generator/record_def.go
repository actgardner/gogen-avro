package generator

import (
	"fmt"
)

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

func (r *recordDefinition) auxStructs() string {
	auxDefs := make(map[string]string)
	imports := make(map[string]string)
	imports["io"] = "import \"io\""
	for _, f := range r.fields {
		f.AuxStructs(auxDefs, imports)
	}
	importStr := concatSortedMap(imports, "\n")
	auxDefStr := concatSortedMap(auxDefs, "\n")
	return importStr + auxDefStr
}

func (r *recordDefinition) fieldSerializers() string {
	serializerMethods := "var err error\n"
	for _, f := range r.fields {
		serializerMethods += fmt.Sprintf("err = %v(r.%v, w)\nif err != nil {return err}\n", f.SerializerMethod(), f.Name())
	}
	return serializerMethods
}

func (r *recordDefinition) structDefinition() string {
	return fmt.Sprintf("type %v struct {\n%v}\n", r.goName(), r.structFields())
}

func (r *recordDefinition) serializerMethod() string {
	return fmt.Sprintf("func (r *%v) Serialize(w io.Writer) error {\n%v\nreturn nil\n}", r.goName(), r.fieldSerializers())
}
