package templates

const EnumTemplate = `
import (
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v7/vm"
)

{{ if ne .Doc "" }}
// {{ .Doc}}
{{ end }}
type {{ .GoType }} int32

const (
{{ range $i, $symbol := .Symbols }}
	{{ $.SymbolName $symbol }} {{ $.GoType }} = {{ $i }}
{{ end }}
)

func (e {{ .GoType  }}) String() string {
	switch e {
{{ range $i, $symbol := .Symbols }}
	case {{ $.SymbolName $symbol }}:
		return {{ printf "%q" $symbol }}
{{ end }}
	}
	return "unknown"
}

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	return vm.WriteInt(int32(r), w)
}

func {{ .FromStringMethod }}(raw string) (r {{ .GoType }}, err error) {
	switch raw {
{{ range $i, $symbol := .Symbols }}
	case {{ printf "%q" $symbol }}:
		return {{ $.SymbolName $symbol }}, nil
{{ end }}
	}

	return -1, fmt.Errorf("invalid value for {{ $.GoType }}: '%s'", raw)
}
`
