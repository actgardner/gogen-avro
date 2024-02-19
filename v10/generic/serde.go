package generic

import (
	"io"

	"github.com/actgardner/gogen-avro/v10/compiler"
	"github.com/actgardner/gogen-avro/v10/schema"
	"github.com/actgardner/gogen-avro/v10/vm"
)

type Codec struct {
	t     schema.AvroType
	deser *vm.Program
}

func NewCodecFromSchema(writer, reader []byte) (*Codec, error) {
	readerType, err := compiler.ParseSchema(reader)
	if err != nil {
		return nil, err
	}

	writerType, err := compiler.ParseSchema(writer)
	if err != nil {
		return nil, err
	}

	return NewCodec(writerType, readerType)
}

func NewCodec(writer, reader schema.AvroType) (*Codec, error) {
	prog, err := compiler.Compile(writer, reader, compiler.GenericMode())
	if err != nil {
		return nil, err
	}

	return &Codec{
		t:     reader,
		deser: prog,
	}, nil
}

func (c *Codec) Deserialize(r io.Reader) (interface{}, error) {
	datum := DatumForType(c.t)
	err := vm.Eval(r, c.deser, datum)
	if err != nil {
		return nil, err
	}
	return datum.Datum(), nil
}
