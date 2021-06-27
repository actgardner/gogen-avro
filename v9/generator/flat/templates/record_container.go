package templates

const RecordContainerTemplate = `
import (
	"io"

	"github.com/actgardner/gogen-avro/v9/container"
	"github.com/actgardner/gogen-avro/v9/vm"
	"github.com/actgardner/gogen-avro/v9/compiler"
)

func {{ .NewWriterMethod }}(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := {{ .ConstructorMethod }}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}

// container reader
type {{ .RecordReaderTypeName }} struct {
	r io.Reader
	p *vm.Program
}

func New{{ .RecordReaderTypeName }}(r io.Reader) (*{{ .RecordReaderTypeName }}, error){
	containerReader, err := container.NewReader(r)
	if err != nil {
		return nil, err
	}

	t := {{ .ConstructorMethod }}
	deser, err := compiler.CompileSchemaBytes([]byte(containerReader.AvroContainerSchema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	return &{{ .RecordReaderTypeName }} {
		r: containerReader,
		p: deser,
	}, nil
}

func (r {{ .RecordReaderTypeName }}) Read() ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}
        err := vm.Eval(r.r, r.p, &t)
	return t, err
}
`
