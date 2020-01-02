package templates

const EnumTemplate = `
import (
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/vm"
)

{{ if ne .Enum.Doc "" }}
// {{ .Enum.Doc}}
{{ end }}  
type {{ .TypeMetadata.GoType }} int32

const (
{{ range $i, $symbol := .Enum.Symbols }}
	{{ $.Enum.SymbolName $symbol }} {{ $.TypeMetadata.GoType }} = {{ $i }}
{{ end }}
)

func (e {{ .TypeMetadata.GoType  }}) String() string {
	switch e {
{{ range $i, $symbol := .Symbols }}
	case {{ $.SymbolName $symbol }}:
		return {{ printf "%q" $symbol }}
{{ end }}
	}
	return "unknown"
}

func {{ .SerializerMethod }}(r {{ .TypeMetadata.GoType }}, w io.Writer) error {
	return vm.WriteInt(int32(r), w)
}

func {{ .FromStringMethod }}(raw string) (r {{ .TypeMetadata.GoType }}, err error) {
	switch raw {
{{ range $i, $symbol := .Symbols }}
	case {{ printf "%q" $symbol }}:
		return {{ $.SymbolName $symbol }}, nil
{{ end }}
	}

	return -1, fmt.Errorf("invalid value for {{ $.TypeMetadata.GoType }}: '%s'", raw)
}
`
