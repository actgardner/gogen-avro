package templates

type RecordContext struct {
	def *model.RecordDefinition
}

func NewRecordContext(def *model.RecordDefinition) *RecordContext {
	return &RecordContext{def}
}

func (r *RecordContext) GoType() string {
	return fmt.Sprintf("*%v", r.def.Name())
}

func (r *RecordContext) ConstructorMethod() string {
	return fmt.Sprintf("New%v()", r.def.Name())
}

func (r *RecordContext) DefaultForField(f *Field) (string, error) {
	return f.Type().DefaultValue(fmt.Sprintf("r.%v", f.GoName()), f.Default())
}

func (r *RecordContext) ConstructableForField(f *Field) string {
	if constructor, ok := getConstructableForType(f.Type()); ok {
		return fmt.Sprintf("r.%v = %v\n", f.GoName(), constructor.ConstructorMethod())
	}
	return ""
}

func (s *RecordContext) WrapperType() string {
	return ""
}

func (r *RecordDefinition) RecordReaderTypeName() string {
	return r.Name() + "Reader"
}

func (r *RecordDefinition) SerializerMethod() string {
	return fmt.Sprintf("write%v", r.Name())
}

func (r *RecordDefinition) NewWriterMethod() string {
	return fmt.Sprintf("New%vWriter", r.Name())
}

func (r *RecordContext) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
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

const RecordTemplate = `
import (
	"io"
	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/compiler"
)

{{ if ne .Doc "" }}
// {{ .Doc}}
{{ end }}  
type {{ .Name }} struct {
{{ range $i, $field := .Fields }}
	{{ if ne $field.Doc "" }}
	// {{ $field.Doc }}
	{{ end }}
	{{ if ne $field.Tags "" }}
		{{ $field.SimpleName }} {{ $field.Type.GoType }} ` + "`{{ $field.Tags }}`" + `
	{{ else }}
		{{ $field.SimpleName }} {{ $field.Type.GoType }}
	{{ end }}
{{ end }}
}

func {{ .ConstructorMethod }} ({{ .GoType}}) {
	return &{{ .Name }}{}
}

func Deserialize{{ .Name }}(r io.Reader) ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err	
	}
	return t, err
}

func Deserialize{{ .Name }}FromSchema(r io.Reader, schema string) ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err	
	}
	return t, err
}

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	var err error
	{{ range $i, $field := .Fields }}
	err = {{ .Type.SerializerMethod }}( r.{{ .GoName }}, w)
	if err != nil {
		return err			
	}
	{{ end }}
	return err
}

func (r {{ .GoType }}) Serialize(w io.Writer) error {
	return {{ .SerializerMethod }}(r, w)
}

func (r {{ .GoType }}) Schema() string {
	return {{ printf "%q" .Schema }}
}

func (r {{ .GoType }}) SchemaName() string {
	return {{ printf "%q" .AvroName.String }}
}

func (_ {{ .GoType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetString(v string) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r {{ .GoType }}) Get(i int) types.Field {
	switch (i) {
	{{ range $i, $field := .Fields }}
	case {{ $i }}:
		{{ if $.ConstructableForField $field | ne "" }}
			{{ $.ConstructableForField $field }}	
		{{ end }}
		{{ if ne $field.Type.WrapperType "" }}
			return (*{{ $field.Type.WrapperType }})(&r.{{ $field.GoName }})
		{{ else }}
			return r.{{ $field.GoName }}
		{{ end }}
	{{ end }}
	}
	panic("Unknown field index")
}

func (r {{ .GoType }}) SetDefault(i int) {
	switch (i) {
	{{ range $i, $field := .Fields }}
        {{ if .HasDefault }}
	case {{ $i }}:
       	 	{{ $.DefaultForField $field }}
		return
	{{ end }}
	{{ end }}
	}
	panic("Unknown field index")
}

func (_ {{ .GoType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) Finalize() { }
`
