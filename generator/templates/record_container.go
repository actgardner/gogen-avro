package templates

const RecordContainerTemplate = `
{{ $metadata := nodeMetadata . }}

import (
	"io"

	"github.com/actgardner/gogen-avro/container"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/compiler"
)

func {{ $metadata.NewWriterMethod }}(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := {{ $metadata.ConstructorMethod }}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}

// container reader
type {{ $metadata.RecordReaderTypeName }} struct {
	r io.Reader
	p *vm.Program
}

func New{{ $metadata.RecordReaderTypeName }}(r io.Reader) (*{{ $metadata.RecordReaderTypeName }}, error){
	containerReader, err := container.NewReader(r)
	if err != nil {
		return nil, err
	}

	t := {{ $metadata.ConstructorMethod }}
	deser, err := compiler.CompileSchemaBytes([]byte(containerReader.AvroContainerSchema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	return &{{ $metadata.RecordReaderTypeName }} {
		r: containerReader,
		p: deser,
	}, nil
}

func (r {{ $metadata.RecordReaderTypeName }}) Read() ({{ $metadata.GoType }}, error) {
	t := {{ $metadata.ConstructorMethod }}
        err := vm.Eval(r.r, r.p, t)
	return t, err
}
`
