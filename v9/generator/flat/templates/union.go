package templates

const UnionTemplate = `
import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v9/vm"
	"github.com/actgardner/gogen-avro/v9/compiler"
	"github.com/actgardner/gogen-avro/v9/vm/types"
)


type {{ .UnionEnumType }} int
const (
{{ range $i, $t := .ItemTypes -}}
	{{ if ne $i $.NullIndex -}}
	{{ $.UnionEnumType }}{{ .Name }} {{ $.UnionEnumType }} = {{ $i }}
        {{ end }}
{{ end -}}
)

type {{ .Name }} struct {
{{ range $i, $t := .ItemTypes -}}
	{{ .Name }} {{ .GoType }}
{{ end -}}
	UnionType {{ $.UnionEnumType }}
}

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	{{ if ne .NullIndex -1 }}
	if r == nil {
		err := vm.WriteLong({{ $.NullIndex }}, w)
		return err
	}
        {{ end }}
	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType{
	{{ range $i, $t := .ItemTypes -}}
	  {{ if ne $i $.NullIndex -}}
	  case {{ $.ItemName $t }}:
		return {{ .SerializerMethod }}(r.{{ .Name }}, w)
          {{ end -}}
	{{ end -}}
	}
	return fmt.Errorf("invalid value for {{ .GoType }}")
}

{{ if ne .NullIndex -1 }}
func {{ .ConstructorMethod }} {{ .GoType }} {
	return &{{ .Name }}{}
} 
{{ else }}
func {{ .ConstructorMethod }} {{ .GoType }} {
	return {{ .Name }}{}
} 
{{ end }}

func (r {{ .GoType }}) Serialize(w io.Writer) error {
	return {{ .SerializerMethod }}(r, w)
}

func Deserialize{{ .Name }}(r io.Reader) ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

{{ if ne .NullIndex -1 }}
	err = vm.Eval(r, deser, t)
{{ else }}
	err = vm.Eval(r, deser, &t)
{{ end }}
	if err != nil {
		return t, err
	}
	return t, err
}

func Deserialize{{ .Name }}FromSchema(r io.Reader, schema string) ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}
	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

{{ if ne .NullIndex -1 }}
	err = vm.Eval(r, deser, t)
{{ else }}
	err = vm.Eval(r, deser, &t)
{{ end }}
	if err != nil {
		return t, err
	}
	return t, err
}

func (r {{ .GoType }}) Schema() string {
	return {{ printf "%q" .Schema }}
}

func (_ {{ .GoType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetString(v string) { panic("Unsupported operation") }
{{ if ne .NullIndex -1 }}
func (r {{ .GoType }}) SetLong(v int64) { 
{{ else }}
func (r *{{ .GoType }}) SetLong(v int64) { 
{{ end }}
	r.UnionType = ({{ .UnionEnumType }})(v)
}
{{ if ne .NullIndex -1 }}
func (r {{ .GoType }}) Get(i int) types.Field {
{{ else }}
func (r *{{ .GoType }}) Get(i int) types.Field {
{{ end }}
	switch (i) {
	{{ range $i, $t := .ItemTypes -}}
	case {{ $i }}:
		{{ if $.ItemConstructor $t | ne "" -}}
		r.{{ .Name }} = {{ $.ItemConstructor $t }}
		{{ end -}}
		{{ if eq .WrapperType "" -}}
		return r.{{ .Name }}
		{{ else -}}
		return &{{ .WrapperType }}{Target: (&r.{{ .Name }})}
		{{ end -}}
	{{ end -}}
	}
	panic("Unknown field index")
}
func (_ {{ .GoType }}) NullField(i int) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetDefault(i int) { panic("Unsupported operation") }
func (_ {{ .GoType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) Finalize()  { }

func (r {{ .GoType }}) MarshalJSON() ([]byte, error) {
	{{ if ne .NullIndex -1 }}
	if r == nil {
		return []byte("null"), nil
	}
	{{ end }}
	switch r.UnionType{
	{{ range $i, $t := .ItemTypes -}}
	{{ if ne $i $.NullIndex -}}
	case {{ $.ItemName $t }}:
		return json.Marshal(map[string]interface{}{"{{ .UnionKey }}": r.{{ .Name }}})
        {{ end -}}
	{{ end -}}
	}
	return nil, fmt.Errorf("invalid value for {{ .GoType }}")
}


{{ if ne .NullIndex -1 }}
func (r {{ .GoType }}) UnmarshalJSON(data []byte) (error) {
{{ else }}
func (r *{{ .GoType }}) UnmarshalJSON(data []byte) (error) {
{{ end }}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}
	{{ range $i, $t := .ItemTypes -}}
	{{ if ne $i $.NullIndex -}}
	if value,  ok := fields["{{ .UnionKey }}"]; ok {
		r.UnionType = {{ $i }}
		return json.Unmarshal([]byte(value), &r.{{ .Name }})
	}
        {{ end -}}
	{{ end -}}
	return fmt.Errorf("invalid value for {{ .GoType }}")
}
`
