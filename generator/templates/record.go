package templates

const RecordTemplate = `
{{ $metadata := nodeMetadata . }}

import (
	"io"
	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/compiler"
)

{{ if ne .Doc "" }}
// {{ .Doc}}
{{ end }}  
type {{ $metadata.Name }} struct {
{{ range $i, $field := .Fields }}
	{{ $fieldMetadata := fieldMetadata $field }}
	{{ $fieldTypeMetadata := nodeMetadata $field.Type }}
	{{ if ne $field.Doc "" }}
	// {{ $field.Doc }}
	{{ end }}
	{{ if ne $field.Tags "" }}
		{{ $fieldMetadata.Name }} {{ $fieldTypeMetadata.GoType }} ` + "`{{ $field.Tags }}`" + `
	{{ else }}
		{{ $fieldMetadata.Name }} {{ $fieldTypeMetadata.GoType }}
	{{ end }}
{{ end }}
}

func {{ $metadata.ConstructorMethod }} ({{ $metadata.GoType}}) {
	return &{{ $metadata.Name }}{}
}

func Deserialize{{ $metadata.Name }}(r io.Reader) ({{ $metadata.GoType }}, error) {
	t := {{ $metadata.ConstructorMethod }}
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

func Deserialize{{ $metadata.Name }}FromSchema(r io.Reader, schema string) ({{ $metadata.GoType }}, error) {
	t := {{ $metadata.ConstructorMethod }}

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

func {{ $metadata.SerializerMethod }}(r {{ $metadata.GoType }}, w io.Writer) error {
	var err error
	{{ range $i, $field := .Fields }}
	{{ $fieldMetadata := fieldMetadata $field }}
	{{ $fieldTypeMetadata := nodeMetadata $field.Type }}

	err = {{ $fieldTypeMetadata.SerializerMethod }}( r.{{ $fieldMetadata.Name }}, w)
	if err != nil {
		return err			
	}
	{{ end }}
	return err
}

func (r {{ $metadata.GoType }}) Serialize(w io.Writer) error {
	return {{ $metadata.SerializerMethod }}(r, w)
}

func (r {{ $metadata.GoType }}) Schema() string {
	return {{ printf "%q" .Schema }}
}

func (r {{ $metadata.GoType }}) SchemaName() string {
	return {{ printf "%q" .AvroName.String }}
}

func (_ {{ $metadata.GoType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetString(v string) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r {{ $metadata.GoType }}) Get(i int) types.Field {
	switch (i) {
	{{ range $i, $field := .Fields }}
	{{ $fieldMetadata := fieldMetadata $field }}
	{{ $fieldTypeMetadata := nodeMetadata $field.Type }}

	case {{ $i }}:
		{{ if $fieldTypeMetadata.ConstructorMethod | ne "" }}
			{{ $fieldMetadata.Name }} = {{ $fieldTypeMetadata.ConstructorMethod  }}
		{{ end }}
		{{ if ne $fieldTypeMetadata.WrapperType "" }}
			return (*{{ $fieldTypeMetadata.WrapperType }})(&r.{{ $fieldMetadata.Name }})
		{{ else }}
			return r.{{ $fieldMetadata.Name }}
		{{ end }}
	{{ end }}
	}
	panic("Unknown field index")
}

func (r {{ $metadata.GoType }}) SetDefault(i int) {
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

func (_ {{ $metadata.GoType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) Finalize() { }
`
