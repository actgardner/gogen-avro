package container

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"

	"github.com/golang/snappy"

	"github.com/actgardner/gogen-avro/container/avro"
	"github.com/actgardner/gogen-avro/schema"
)

// Reader is a low-level primitive for reading the OCF framing of a file.
// Generally you can create a Reader using the `New<RecordType>Reader` method generate for every record type.
type Reader struct {
	codec            Codec
	reader           io.Reader
	compressedReader io.Reader
	schemaBytes      []byte
	schema           schema.AvroType
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
	log("Got OCF schema from header: %v", string(schemaBytes))

	codec, ok := header.Meta["avro.codec"]
	if !ok {
		return nil, fmt.Errorf("Expected avro.codec in header, not specified in metadata map - %v", header.Meta)
	}
	log("Got OCF codec from header: %v", string(codec))

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
		log("OCF reader opening new block")
		if err := r.openBlock(); err != nil {
			return 0, err
		}
	}

	for {
		n, err := r.compressedReader.Read(b)
		log("OCF container read: %v %v", n, err)
		if n > 0 {
			return n, nil
		}

		if err == io.EOF {
			log("OCF EOF, opening new block")
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

	log("OCF block size: %v", len(header.RecordBytes))
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
