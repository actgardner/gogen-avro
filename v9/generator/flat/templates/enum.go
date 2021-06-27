package templates

const EnumTemplate = `
import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v9/vm/types"
	"github.com/actgardner/gogen-avro/v9/vm"
)

{{ if ne .Doc "" }}// {{ .Doc}}{{ end }}  
type {{ .GoType }} int32

const (
{{ range $i, $symbol := .Symbols -}}
	{{ $.SymbolName $symbol }} {{ $.GoType }} = {{ $i }}
{{ end -}}
)

func (e {{ .GoType  }}) String() string {
	switch e {
	{{ range $i, $symbol := .Symbols -}}
	case {{ $.SymbolName $symbol }}:
		return {{ printf "%q" $symbol }}
	{{ end -}}
	}
	return "unknown"
}

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	return vm.WriteInt(int32(r), w)
}

func {{ .FromStringMethod }}(raw string) (r {{ .GoType }}, err error) {
	switch raw {
{{ range $i, $symbol := .Symbols -}}
	case {{ printf "%q" $symbol }}:
		return {{ $.SymbolName $symbol }}, nil
{{ end -}}
	}

	return -1, fmt.Errorf("invalid value for {{ $.GoType }}: '%s'", raw)
}

func (b {{ .GoType }}) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *{{ .GoType }}) UnmarshalJSON(data []byte) (error) {
	var stringVal string
	err := json.Unmarshal(data, &stringVal)
	if err != nil {
		return err
	}
	val, err := {{ .FromStringMethod }}(stringVal)
	*b = val
	return err
}

type {{ .WrapperType }} struct {
	Target *{{ .GoType }}
}

func (b *{{ .WrapperType }}) SetBoolean(v bool) {
	panic("Unable to assign boolean to int field")
}

func (b *{{ .WrapperType }}) SetInt(v int32) {
	*(b.Target) = {{ .GoType }}(v)
}

func (b *{{ .WrapperType }}) SetLong(v int64) {
	panic("Unable to assign long to int field")
}

func (b *{{ .WrapperType }}) SetFloat(v float32) {
	panic("Unable to assign float to int field")
}

func (b *{{ .WrapperType }}) SetUnionElem(v int64) {
	panic("Unable to assign union elem to int field")
}

func (b *{{ .WrapperType }}) SetDouble(v float64) {
	panic("Unable to assign double to int field")
}

func (b *{{ .WrapperType }}) SetBytes(v []byte) {
	panic("Unable to assign bytes to int field")
}

func (b *{{ .WrapperType }}) SetString(v string) {
	panic("Unable to assign string to int field")
}

func (b *{{ .WrapperType }}) Get(i int) types.Field {
	panic("Unable to get field from int field")
}

func (b *{{ .WrapperType }}) SetDefault(i int) {
	panic("Unable to set default on int field")
}

func (b *{{ .WrapperType }}) AppendMap(key string) types.Field {
	panic("Unable to append map key to from int field")
}

func (b *{{ .WrapperType }}) AppendArray() types.Field {
	panic("Unable to append array element to from int field")
}


func (b *{{ .WrapperType }}) NullField(int) {
	panic("Unable to null field in int field")
}

func (b *{{ .WrapperType }}) Finalize() {}
`
