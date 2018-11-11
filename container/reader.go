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

/**
 * Reader provides an experimental Object Container File reader to load Avro container files. You can
 * give this an OCF and it'll transparently decode a block at a time, so you can pass this to your
 * generated deserializers. See `test/primitive/container_test.go` for an example.
 *
 * Note: This is experimental and the interface may change or be deprecated. Right now gogen-avro
 *       only supports deserializing with the exact schema you used to serialize - there is no
 *       support for adding/removing fields or changing types.
 */
type Reader struct {
	codec            Codec
	reader           io.Reader
	compressedReader io.Reader
	schemaBytes      []byte
	schema           types.AvroType
	sync             avro.Sync
}

func NewReader(r io.Reader) (*Reader, error) {
	header, err := avro.DeserializeAvroContainerHeader(r)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(header.Magic[:], []byte{'o', 'b', 'j', 1}) {
		return nil, fmt.Errorf("Unexpected magic in header - %v", header.Magic)
	}

	schemaBytes, ok := header.Meta["avro.schema"]
	if !ok {
		return nil, fmt.Errorf("Expected avro.schema in header, not specified in metadata map - %v", header.Meta)
	}

	codec, ok := header.Meta["avro.codec"]
	if !ok {
		return nil, fmt.Errorf("Expected avro.codec in header, not specified in metadata map - %v", header.Meta)
	}

	return &Reader{
		codec:            Codec(codec),
		reader:           r,
		schemaBytes:      schemaBytes,
		compressedReader: nil,
		schema:           nil,
		sync:             header.Sync,
	}, nil
}

func (r *Reader) AvroContainerSchema() []byte {
	return r.schemaBytes
}

func (r *Reader) Read(b []byte) (n int, err error) {
	if r.compressedReader == nil {
		if err := r.openBlock(); err != nil {
			return 0, err
		}
	}

	for {
		n, err := r.compressedReader.Read(b)
		if n > 0 {
			return n, nil
		}

		if err == io.EOF {
			if err := r.openBlock(); err != nil {
				return 0, err
			}
			continue
		}
		return n, err
	}
}

func (r *Reader) openBlock() error {
	header, err := avro.DeserializeAvroContainerBlock(r.reader)
	if err != nil {
		return err
	}

	if header.Sync != r.sync {
		return fmt.Errorf("Unexpected sync marker %q, expected %q", header.Sync, r.sync)
	}

	blockBuffer := bytes.NewBuffer(header.RecordBytes)

	switch r.codec {
	case Null:
		r.compressedReader = blockBuffer
		break
	case Deflate:
		r.compressedReader = flate.NewReader(blockBuffer)
		break
	case Snappy:
		// TODO: Check the last 4 bytes are the big-endian CRC of the compressed Snappy block
		dst := make([]byte, 0, 0)
		dst, err := snappy.Decode(nil, header.RecordBytes[:len(header.RecordBytes)-4])
		if err != nil {
			return err
		}
		r.compressedReader = bytes.NewBuffer(dst)
		break
	default:
		return fmt.Errorf("Unexpected codec %q", r.codec)
	}

	return nil
}
