package templates

const UnionTemplate = `
import (
	"errors"
	"io"

	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
)


type {{ .UnionEnumType }} int

const ({{ range $i, $t := .AvroTypes }}{{ if $.OptionalIndex | ne $i }}
	{{ $.UnionEnumType }}{{ .Name }} {{ $.UnionEnumType }} = {{ $i }}
{{ end }}{{ end }})

type {{ .Name }} struct { {{ range $i, $t := .ItemTypes }}
	{{ .Name }} {{ .GoType }}
{{ end }}
	UnionType {{ $.UnionEnumType }}
}

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error { {{ if .IsOptional }}
	if r == nil {
		return vm.WriteLong(int64({{ .OptionalIndex }}), w)
	} {{ end }}
	if err := vm.WriteLong(int64(r.UnionType), w); err != nil {
		return err
	}
	switch r.UnionType{ {{ range $i, $t := .ItemTypes }}
	case {{ $.ItemName $t }}:
		return {{ .SerializerMethod }}(r.{{ .Name }}, w){{ end }}
	}
	return errors.New("invalid value for {{ .GoType }}")
}

func (_ *{{ .Name }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *{{ .Name }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ *{{ .Name }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *{{ .Name }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *{{ .Name }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *{{ .Name }}) SetString(v string) { panic("Unsupported operation") }

func (r *{{ .Name }}) SetLong(v int64) { 
	r.UnionType = ({{ .UnionEnumType }})(v)
}

func (r *{{ .Name }}) Get(i int) types.Field {
	switch (i) { {{ range $i, $t := .AvroTypes }}{{ if $.OptionalIndex | ne $i }}
	case {{ $i }}:
		{{ if $.ItemConstructor $t | ne "" }}
		r.{{ .Name }} = {{ $.ItemConstructor $t }}
		{{ end }}
		{{ if eq .WrapperType "" }}
		return &r.{{ .Name }}
		{{ else }}
		return (*{{ .WrapperType }})(&r.{{ .Name }})
		{{ end }}
	{{ end }}{{ end }}
	}
	panic("Unknown field index")
}

func (r *{{ .Name }}) Clear(i int) { panic("Unsupported operation") }
func (_ *{{ .Name }}) SetDefault(i int) { panic("Unsupported operation") }
func (_ *{{ .Name }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *{{ .Name }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *{{ .Name }}) Finalize()  { }
`
