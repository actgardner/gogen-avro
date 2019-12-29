package templates

const FixedTemplate = `
import (
	"io"
	"github.com/actgardner/gogen-avro/vm/types"
)

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	_, err := w.Write(r[:])
	return err
}

type {{ .GoType }} {{ .WrapperType }}
type {{ .WrapperType }} [{{ .SizeBytes }}]byte

func (_ *{{ .WrapperType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (r *{{ .WrapperType }}) SetBytes(v []byte) { 
	copy((*r)[:], v)
}
func (_ *{{ .WrapperType }}) SetString(v string) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) Finalize() { }
func (_ *{{ .WrapperType }}) SetDefault(i int) { panic("Unsupported operation") }

`
