package templates

const EnumTemplate = `
import (
	"io"

	"github.com/actgardner/gogen-avro/vm"
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
`
