package container

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"

	"github.com/golang/snappy"

	"github.com/actgardner/gogen-avro/container/avro"
	"github.com/actgardner/gogen-avro/types"
)

type Reader struct {
	codec            Codec
	reader           io.Reader
	compressedReader io.Reader
	schema           types.AvroType
	sync             avro.Sync
}

func NewReader(r io.Reader) (*Reader, error) {
	header, err := avro.DeserializeAvroContainerHeader(r)
	if err != nil {
		return nil, err
	}

	if header.Magic != [4]byte{'o', 'b', 'j', '1'} {
		return nil, fmt.Errorf("Unexpected magic in header - %v", header.Magic)
	}

	schemaString, ok := header.Meta["avro.schema"]
	if !ok {
		return nil, fmt.Errorf("Expected avro.schema in header, not specified in metadata map - %v", header.Meta)
	}

	codec, ok := header.Meta["avro.codec"]
	if !ok {
		return nil, fmt.Errorf("Expected avro.codec in header, not specified in metadata map - %v", header.Meta)
	}

	namespace := types.NewNamespace(false)
	schema, err := namespace.TypeForSchema(schemaString)
	if err != nil {
		return nil, err
	}

	return &Reader{
		codec:            Codec(codec),
		reader:           r,
		compressedReader: nil,
		schema:           schema,
		sync:             header.Sync,
	}, nil
}

func (r *Reader) Read(b []byte) (n int, err error) {
	if r.compressedReader == nil {
		if err := r.openBlock(); err != nil {
			return 0, err
		}
	}

	for {
		n, err := r.compressedReader.Read(b)
		if err == io.EOF {
			if err := r.openBlock(); err != nil {
				return 0, err
			}
		}
		return n, err
	}
}

func (r *Reader) openBlock() error {
	header, err := avro.DeserializeAvroContainersBlock(r.reader)
	if err != nil {
		return err
	}

	if header.Sync != r.sync {
		return fmt.Errorf("Unexpected sync marker %q, expected %q", header.Sync, r.sync)
	}

	blockBuffer := bytes.Buffer(header.RecordBytes)

	switch r.codec {
	case Null:
		r.compressedReader = blockBuffer
		break
	case Deflate:
		r.compressedReader = flate.NewReader(blockBuffer)
		break
	case Snappy:
		r.compressedReader = snappy.NewReader(blockBuffer)
		break
	default:
		return fmt.Errorf("Unexpected codec %q", r.codec)
	}

	return nil
}
