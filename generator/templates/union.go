package templates

const UnionTemplate = `
{{ $metadata := nodeMetadata . }}

import (
	"io"
	"fmt"

	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
)


type {{ $metadata.UnionEnumType }} int
const (
{{ range $i, $t := .ItemTypes }}
	 {{ $metadata.UnionEnumType }}{{ .Name }} {{ $metadata.UnionEnumType }} = {{ $i }}
{{ end }}
)

type {{ .Name }} struct {
{{ range $i, $t := .ItemTypes }}
	{{ .Name }} {{ $metadata.GoType }}
{{ end }}
	UnionType {{ $metadata.UnionEnumType }}
}

func {{ .SerializerMethod }}(r {{ $metadata.GoType }}, w io.Writer) error {
	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType{
	{{ range $i, $t := .ItemTypes }}
	case {{ $.ItemName $t }}:
		return {{ .SerializerMethod }}(r.{{ .Name }}, w)
        {{ end }}
	}
	return fmt.Errorf("invalid value for {{ $metadata.GoType }}")
}

func {{ $metadata.ConstructorMethod }} {{ $metadata.GoType }} {
	return &{{ $metadata.Name }}{}
}

func (_ {{ $metadata.GoType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) SetString(v string) { panic("Unsupported operation") }
func (r {{ $metadata.GoType }}) SetLong(v int64) { 
	r.UnionType = ({{ $metadata.UnionEnumType }})(v)
}
func (r {{ $metadata.GoType }}) Get(i int) types.Field {
	switch (i) {
	{{ range $i, $t := .ItemTypes }}
	case {{ $i }}:
		{{ if $.ItemConstructor $t | ne "" }}
		r.{{ .Name }} = {{ $.ItemConstructor $t }}
		{{ end }}
		{{ if eq .WrapperType "" }}
		return r.{{ .Name }}
		{{ else }}
		return (*{{ .WrapperType }})(&r.{{ .Name }})
		{{ end }}
	{{ end }}
	}
	panic("Unknown field index")
}
func (_ {{ $metadata.GoType }}) SetDefault(i int) { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ {{ $metadata.GoType }}) Finalize()  { }
`
