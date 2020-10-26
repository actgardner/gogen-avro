package templates

const FixedTemplate = `
import (
	"io"
	"encoding/json"

	"github.com/actgardner/gogen-avro/v7/vm/types"
)

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	_, err := w.Write(r[:])
	return err
}

type {{ .WrapperType }} struct {
	Target *{{ .GoType }}
}

type {{ .GoType }} [{{ .SizeBytes }}]byte

func (b *{{ .GoType }}) UnmarshalJSON(data []byte) (error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	copy((*b)[:], []byte(s))
	return nil
}

func (b {{ .GoType }}) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(string(b[:]))
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (_ *{{ .WrapperType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (r *{{ .WrapperType }}) SetBytes(v []byte) { 
	copy((*r.Target)[:], v)
}
func (_ *{{ .WrapperType }}) SetString(v string) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) NullField(int) { panic("Unsupported operation") }
func (_ *{{ .WrapperType }}) Finalize() { }
func (_ *{{ .WrapperType }}) SetDefault(i int) { panic("Unsupported operation") }

`
