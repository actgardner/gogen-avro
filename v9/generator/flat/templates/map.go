package templates

const MapTemplate = `
import (
	"io"
	"github.com/actgardner/gogen-avro/v9/vm/types"
	"github.com/actgardner/gogen-avro/v9/vm"
)

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	err := vm.WriteLong(int64(len(r)), w)
	if err != nil || len(r) == 0 {
		return err
	}
	for k, e := range r {
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

type {{ .WrapperType }} struct {
	Target *map[string]{{ .ItemType.GoType }}
	keys []string
	values []{{ .ItemType.GoType }}
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
func (_ *{{ .WrapperType }}) SetDefault(i int) { panic("Unsupported operation") }

func (r *{{ .WrapperType }}) NullField(_ int) { 
	{{ if isNullable .ItemType -}}
		r.values[len(r.values)-1] = nil		
	{{ else -}}
		panic("Unsupported operation")
	{{ end -}}
}

func (r *{{ .WrapperType }}) Finalize() { 
	for i := range r.keys {
		(*r.Target)[r.keys[i]] = r.values[i]
	}
}

func (r *{{ .WrapperType }}) AppendMap(key string) types.Field { 
	r.keys = append(r.keys, key)
	var v {{ .ItemType.GoType }}
	{{ if ne .ItemConstructable "" -}}
		{{ .ItemConstructable }}
	{{ end -}}
	r.values = append(r.values, v)
	{{ if .ItemType.WrapperType | ne "" -}}
	return &{{ .ItemType.WrapperType }}{Target: &r.values[len(r.values)-1]}
	{{ else -}}
	return r.values[len(r.values)-1]
	{{ end -}}
}

func (_ *{{ .WrapperType }}) AppendArray() types.Field { panic("Unsupported operation") }

`
