package templates

const ArrayTemplate = `
import (
	"io"

	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
)

func {{ .NodeWrapper.SerializerMethod }}(r {{ .TypeMetadata.GoType }}, w io.Writer) error {
	err := vm.WriteLong(int64(len(r)),w)
	if err != nil || len(r) == 0 {
		return err
	}
	for _, e := range r {
		err = {{ .ItemType.NodeWrapper.SerializerMethod }}(e, w)
		if err != nil {
			return err
		}
	}
	return vm.WriteLong(0,w)
}



type {{ .TypeMetadata.WrapperType }} {{ .TypeMetadata.GoType }}

func (_ *{{ .TypeMetadata.WrapperType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *{{ .TypeMetadata.WrapperType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ *{{ .TypeMetadata.WrapperType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ *{{ .TypeMetadata.WrapperType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *{{ .TypeMetadata.WrapperType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *{{ .TypeMetadata.WrapperType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *{{ .TypeMetadata.WrapperType }}) SetString(v string) { panic("Unsupported operation") }
func (_ *{{ .TypeMetadata.WrapperType }}) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *{{ .TypeMetadata.WrapperType }}) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *{{ .TypeMetadata.WrapperType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *{{ .TypeMetadata.WrapperType }}) Finalize() { }
func (_ *{{ .TypeMetadata.WrapperType }}) SetDefault(i int) { panic("Unsupported operation") }
func (r *{{ .TypeMetadata.WrapperType }}) AppendArray() types.Field {
	var v {{ (index .Children 0).TypeMetadata.GoType }}
	{{ if (index .Children 0).TypeMetadata.ConstructorMethod }}
	{{ (index .Children 0).TypeMetadata.ConstructorMethod }}
 	{{ end }}
	*r = append(*r, v)
        {{ if (index .Children 0).TypeMetadata.WrapperType }} 
        return (*{{ (index .Children 0).TypeMetadata.WrapperType }})(&(*r)[len(*r)-1])
        {{ else }}
        return (*r)[len(*r)-1]
        {{ end }}
}
`
