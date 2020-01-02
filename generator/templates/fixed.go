package templates

const FixedTemplate = `
{{ $metadata := nodeMetadata . }}
import (
	"io"
	"github.com/actgardner/gogen-avro/vm/types"
)

func {{ $metadata.SerializerMethod }}(r {{ $metadata.GoType }}, w io.Writer) error {
	_, err := w.Write(r[:])
	return err
}

type {{ $metadata.GoType }} {{ $metadata.WrapperType }}
type {{ $metadata.WrapperType }} [{{ .SizeBytes }}]byte

func (_ *{{ $metadata.WrapperType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *{{ $metadata.WrapperType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ *{{ $metadata.WrapperType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ *{{ $metadata.WrapperType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *{{ $metadata.WrapperType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (r *{{ $metadata.WrapperType }}) SetBytes(v []byte) { 
	copy((*r)[:], v)
}
func (_ *{{ $metadata.WrapperType }}) SetString(v string) { panic("Unsupported operation") }
func (_ *{{ $metadata.WrapperType }}) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *{{ $metadata.WrapperType }}) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *{{ $metadata.WrapperType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *{{ $metadata.WrapperType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *{{ $metadata.WrapperType }}) Finalize() { }
func (_ *{{ $metadata.WrapperType }}) SetDefault(i int) { panic("Unsupported operation") }
`
