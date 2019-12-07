package templates

type MapContext struct {
	Field *schema.MapField
}

func (s *MapContext) GoType() string {
	return fmt.Sprintf("*%v", s.Name())
}

func (s *MapContext) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Field.Name())
}

func (s *MapContext) ConstructorMethod() string {
	return fmt.Sprintf("make(%v)", s.Field.Name())
}

func (s *MapContext) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	items, ok := rvalue.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Expected map as default for %v, got %v", lvalue, rvalue)
	}
	setters := ""

	for k, v := range items {
		setter, err := s.itemType.DefaultValue(fmt.Sprintf("%v[%q]", lvalue, k), v)
		if err != nil {
			return "", err
		}
		setters += setter + "\n"
	}
	return setters, nil
}

func (a *MapContext) Template() string {
	return mapTemplate
}

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

type {{ .Name }} struct {
	keys []string
	values []{{ .ItemType.GoType }}
	M map[string]{{ .ItemType.GoType }}
}

func New{{ .Name }}() *{{ .Name }}{
	return &{{ .Name }} {
		keys: make([]string, 0),
		values: make([]{{ .ItemType.GoType }}, 0),
		M: make(map[string]{{ .ItemType.GoType }}),
	}
}

func (_ {{ .GoType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetString(v string) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetUnionElem(v int64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) Get(i int) types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetDefault(i int) { panic("Unsupported operation") }
func (r {{ .GoType }}) Finalize() { 
	for i := range r.keys {
		r.M[r.keys[i]] = r.values[i]
	}
	r.keys = nil
	r.values = nil
}

func (r {{ .GoType }}) AppendMap(key string) types.Field { 
	r.keys = append(r.keys, key)
	var v {{ .ItemType.GoType }}
	{{ if ne .ItemConstructable "" }}
		{{ .ItemConstructable }}
	{{ end }}
	r.values = append(r.values, v)
	{{ if .ItemType.WrapperType | ne "" }}
	return (*{{ .ItemType.WrapperType }})(&r.values[len(r.values)-1])
	{{ else }}
	return r.values[len(r.values)-1]
	{{ end }}
}

func (_ {{ .GoType }}) AppendArray() types.Field { panic("Unsupported operation") }

`
