package templates

const MapTemplate = `
import (
	"io"
	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
)

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	err := vm.WriteLong(int64(len(r.M)), w)
	if err != nil || len(r.M) == 0 {
		return err
	}
	for k, e := range r.M {
		err = vm.WriteString(k, w)
		if err != nil {
			return err
		}
		err = {{ .ItemType.SerializerMethod }}(e, w)
		if err != nil {
			return err
		}
	}
	return vm.WriteLong(0, w)
}

type {{ .GoType }} struct {
	keys []string
	values []{{ .ItemType.GoType }}
	M map[string]{{ .ItemType.GoType }}
}

func New{{ .GoType }}() {{ .GoType }}{
	return {{ .GoType }}{ M: make(map[string]{{ .ItemType.GoType }}) }
}

func (_ *{{ .GoType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetString(v string) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ *{{ .GoType }}) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *{{ .GoType }}) SetDefault(i int) { panic("Unsupported operation") }
func (r *{{ .GoType }}) Finalize() {
	for i := range r.keys {
		r.M[r.keys[i]] = r.values[i]
	}
	r.keys = nil
	r.values = nil
}

func (r *{{ .GoType }}) AppendMap(key string) types.Field { 
	r.keys = append(r.keys, key)
	var v {{ .ItemType.GoType }}
	{{ if ne .ItemConstructable "" }}
		{{ .ItemConstructable }}
	{{ end }}
	r.values = append(r.values, v)
	{{ if .ItemType.WrapperType | ne "" }}
	return (*{{ .ItemType.WrapperType }})(&r.values[len(r.values)-1])
	{{ else }}
	return {{ if not .ItemType.IsOptional }}&{{ end }}r.values[len(r.values)-1]
	{{ end }}
}

func (r *{{ .GoType }}) Clear(i int) { {{ if .ItemType.IsOptional }}
	pos := len(r.keys) - 1
	switch {
	case pos < 0:
		panic("Map already empty")
	default:
		r.values[pos] = nil
	}{{ else }}
	panic("Non-optional map item")
	{{ end }}
}

func (_ *{{ .GoType }}) AppendArray() types.Field { panic("Unsupported operation") }

`
