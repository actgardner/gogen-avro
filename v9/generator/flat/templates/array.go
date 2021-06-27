package templates

const ArrayTemplate = `
import (
	"io"

	"github.com/actgardner/gogen-avro/v9/vm/types"
	"github.com/actgardner/gogen-avro/v9/vm"
)

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	err := vm.WriteLong(int64(len(r)),w)
	if err != nil || len(r) == 0 {
		return err
	}
	for _, e := range r {
		err = {{ .ItemType.SerializerMethod }}(e, w)
		if err != nil {
			return err
		}
	}
	return vm.WriteLong(0,w)
}



type {{ .WrapperType }} struct {
	Target *{{ .GoType }}
}

func (_ *{{ .WrapperType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetString(v string) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) Finalize() { }
func (_ *{{ .WrapperType }}) SetDefault(i int) { panic("Unsupported operation") }
func (r *{{ .WrapperType }}) NullField(i int) { 
	{{ if isNullable .ItemType -}}
		(*r.Target)[len(*r.Target)-1] = nil		
	{{ else -}}
		panic("Unsupported operation")
	{{ end -}}
}

func (r *{{ .WrapperType }}) AppendArray() types.Field {
	var v {{ .ItemType.GoType }}
	{{ if .ItemConstructable -}}
	{{ .ItemConstructable }}
 	{{- end }}
	*r.Target = append(*r.Target, v)
        {{ if .ItemType.WrapperType -}} 
        return &{{ .ItemType.WrapperType }}{Target: &(*r.Target)[len(*r.Target)-1]}
        {{- else }}
        return (*r.Target)[len(*r.Target)-1]
        {{- end }}
}
`
