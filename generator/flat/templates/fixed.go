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

type {{ .GoType }} [{{ .SizeBytes }}]byte

func (_ *{{ .GoType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (r *{{ .GoType }}) SetBytes(v []byte) { 
	copy((*r)[:], v)
}
func (_ *{{ .GoType }}) SetString(v string) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *{{ .GoType }}) Clear(i int) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *{{ .GoType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *{{ .GoType }}) Finalize() { }
func (_ *{{ .GoType }}) SetDefault(i int) { panic("Unsupported operation") }

`
