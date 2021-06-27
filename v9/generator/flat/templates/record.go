package templates

const RecordTemplate = `
import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v9/vm/types"
	"github.com/actgardner/gogen-avro/v9/vm"
	"github.com/actgardner/gogen-avro/v9/compiler"
)

var _ = fmt.Printf

{{ if ne .Doc "" }}// {{ .Doc}}{{ end }}
type {{ .Name }} struct {
{{ range $i, $field := .Fields -}}
	{{ if ne $field.Doc "" }}// {{ $field.Doc }}{{ end }}
	{{ if ne $field.Tags "" -}}
		{{ $field.GoName }} {{ $field.Type.GoType }} ` + "`{{ $field.Tags }}`" + `
	{{ else -}}
		{{ $field.GoName }} {{ $field.Type.GoType }}
	{{ end -}}
{{ end }}
}

const {{ .Name }}AvroCRC64Fingerprint = {{ definitionFingerprint . }}

func {{ .ConstructorMethod }} ({{ .GoType}}) {
	r := {{ .Name }}{}
	{{ range $i, $field := .Fields -}}
	{{ if $.ConstructableForField $field | ne "" -}}
		{{ $.ConstructableForField $field }}
	{{ end -}}
        {{ if .HasDefault -}}
       	 	{{ $.DefaultForField $field }}
	{{ end -}}
	{{ end -}}
	return r
}

func Deserialize{{ .Name }}(r io.Reader) ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func Deserialize{{ .Name }}FromSchema(r io.Reader, schema string) ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	var err error
	{{ range $i, $field := .Fields -}}
	err = {{ .Type.SerializerMethod }}( r.{{ .GoName }}, w)
	if err != nil {
		return err
	}
	{{ end -}}
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

func (r *{{ .GoType }}) Get(i int) types.Field {
	switch (i) {
	{{ range $i, $field := .Fields -}}
	case {{ $i }}:
		{{ if $.ConstructableForField $field | ne "" -}}
			{{ $.ConstructableForField $field }}
		{{ end -}}
		{{ if ne $field.Type.WrapperType "" -}}
			return &{{ $field.Type.WrapperType }}{Target: &r.{{ $field.GoName }}}
		{{ else -}}
			return r.{{ $field.GoName }}
		{{ end -}}
	{{ end -}}
	}
	panic("Unknown field index")
}

func (r *{{ .GoType }}) SetDefault(i int) {
	switch (i) {
	{{ range $i, $field := .Fields -}}
        {{ if .HasDefault -}}
	case {{ $i }}:
       	 	{{ $.DefaultForField $field }}
		return
	{{ end -}}
	{{ end -}}
	}
	panic("Unknown field index")
}

func (r *{{ .GoType }}) NullField(i int) { 
	switch (i) {
	{{ range $i, $field := .Fields -}}
        {{ if isNullable $field.Type -}}
	case {{ $i }}:
		r.{{ $field.GoName }} = nil
		return
	{{ end -}}
	{{ end -}}
	}
	panic("Not a nullable field index")
}

func (_ {{ .GoType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) Finalize() { }


func (_ {{ .GoType}}) AvroCRC64Fingerprint() []byte {
  return []byte({{ .Name }}AvroCRC64Fingerprint)
}

func (r {{ .GoType }}) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	{{ range $i, $field := .Fields -}}
	output[{{ printf "%q" $field.Name }}], err = json.Marshal(r.{{ $field.GoName}})
	if err != nil {
		return nil, err
	}
	{{ end -}}
	return json.Marshal(output)	
}

func (r *{{ .GoType }}) UnmarshalJSON(data []byte) (error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	{{ range $i, $field := .Fields -}}
	val = func() json.RawMessage {
		if v, ok := fields[{{ printf "%q" $field.Name }}]; ok {
			return v
		}
		{{ range $j, $alias := $field.Aliases -}}
		if v, ok := fields[{{ printf "%q" $alias }}]; ok {
			return v
		}
		{{ end -}}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.{{ $field.GoName}}); err != nil {
			return err
		}
	} else {
        	{{ if .HasDefault -}}
		{{ if $.ConstructableForField $field | ne "" -}}
		{{ $.ConstructableForField $field }}
		{{ end -}}
       	 	{{ $.DefaultForField $field }}
		{{ else -}}
		return fmt.Errorf("no value specified for {{ $field.Name }}")
		{{ end -}}
	}
	{{ end -}}
	return nil
}
`
