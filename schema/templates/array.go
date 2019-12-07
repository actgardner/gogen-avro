package templates

type ArrayContext struct {
	Field *schema.ArrayField
}

func (s *ArrayContext) GoType() string {
	return fmt.Sprintf("[]%v", ContextForField(s.Field.ItemType()).GoType())
}

func (s *ArrayContext) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.Field.Name())
}

func (s *ArrayContext) ConstructorMethod() string {
	return fmt.Sprintf("make(%v, 0)", s.GoType())
}

func (s *ArrayContext) WrapperType() string {
	return fmt.Sprintf("%vWrapper", s.Field().Name())
}

func (s *ArrayContext) ItemConstructable() string {
	itemConstructor := ContextForField(s.Field.ItemType())
	if itemConstructor != "" {
		return fmt.Sprintf("v = %v\n")
	}
	return ""
}

func (s *ArrayContext) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	items, ok := rvalue.([]interface{})
	if !ok {
		return "", fmt.Errorf("Expected array as default for %v, got %v", lvalue, rvalue)
	}

	setters := fmt.Sprintf("%v = make(%v,%v)\n", lvalue, s.GoType(), len(items))
	itemConstructor := ContextForField(s.Field.ItemType())
	for i, item := range items {
		if itemConstructor != "" {
			setters += fmt.Sprintf("%v[%v] = %v\n", lvalue, i, itemConstructor)
		}

		setter, err := s.itemType.DefaultValue(fmt.Sprintf("%v[%v]", lvalue, i), item)
		if err != nil {
			return "", err
		}

		setters += setter + "\n"
	}
	return setters, nil
}

func (a *ArrayContext) Template() string {
	return arrayTemplate
}

const arrayTemplate = `
import (
	"io"

	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
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



type {{ .WrapperType }} {{ .GoType }}

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
func (r *{{ .WrapperType }}) AppendArray() types.Field {
	var v {{ .ItemType.GoType }}
	{{ if .ItemConstructable }}
	{{ .ItemConstructable }}
 	{{ end }}
	*r = append(*r, v)
        {{ if .ItemType.WrapperType }} 
        return (*{{ .ItemType.WrapperType }})(&(*r)[len(*r)-1])
        {{ else }}
        return (*r)[len(*r)-1]
        {{ end }}
}
`
