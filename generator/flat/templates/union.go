package templates

const UnionTemplate = `
import (
	"io"
	"fmt"

	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
)


type {{ .UnionEnumType }} int
const (
{{ range $i, $t := .ItemTypes -}}
	 {{ $.UnionEnumType }}{{ .Name }} {{ $.UnionEnumType }} = {{ $i }}
{{ end -}}
)

type {{ .Name }} struct {
{{ range $i, $t := .ItemTypes -}}
	{{ .Name }} {{ .GoType }}
{{ end -}}
	UnionType {{ $.UnionEnumType }}
}

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType{
	{{ range $i, $t := .ItemTypes -}}
	case {{ $.ItemName $t }}:
		return {{ .SerializerMethod }}(r.{{ .Name }}, w)
        {{ end -}}
	}
	return fmt.Errorf("invalid value for {{ .GoType }}")
}

func {{ .ConstructorMethod }} {{ .GoType }} {
	return &{{ .Name }}{}
}

func (_ {{ .GoType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetString(v string) { panic("Unsupported operation") }
func (r {{ .GoType }}) SetLong(v int64) { 
	r.UnionType = ({{ .UnionEnumType }})(v)
}
func (r {{ .GoType }}) Get(i int) types.Field {
	switch (i) {
	{{ range $i, $t := .ItemTypes -}}
	case {{ $i }}:
		{{ if $.ItemConstructor $t | ne "" -}}
		r.{{ .Name }} = {{ $.ItemConstructor $t }}
		{{ end -}}
		{{ if eq .WrapperType "" -}}
		return r.{{ .Name }}
		{{ else -}}
		return (*{{ .WrapperType }})(&r.{{ .Name }})
		{{ end -}}
	{{ end -}}
	}
	panic("Unknown field index")
}
func (_ {{ .GoType }}) SetDefault(i int) { panic("Unsupported operation") }
func (_ {{ .GoType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) Finalize()  { }
`
