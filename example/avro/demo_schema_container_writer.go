package avro

import (
	"bytes"
	"compress/flate"
	"github.com/golang/snappy"
	"io"
)

const (
	DemoSchemaSchema = "{\n\t\"type\": \"record\",\n\t\"name\": \"DemoSchema\",\n\t\"fields\": [\n\t\t{\"name\": \"IntField\", \"type\": \"int\"},\n\t\t{\"name\": \"DoubleField\", \"type\": \"double\"},\n\t\t{\"name\": \"StringField\", \"type\": \"string\"},\n\t\t{\"name\": \"BoolField\", \"type\": \"boolean\"},\n\t\t{\"name\": \"BytesField\", \"type\": \"bytes\"}\n\t]\n\n}\n"
)

type DemoSchemaContainerWriter struct {
	writer          io.Writer
	syncMarker      [16]byte
	codec           Codec
	recordsPerBlock int64

	blockBuffer      *bytes.Buffer
	compressedWriter io.Writer
	nextBlockRecords int64
}

func NewDemoSchemaContainerWriter(writer io.Writer, codec Codec, recordsPerBlock int64) (*DemoSchemaContainerWriter, error) {
	blockBytes := make([]byte, 0)
	blockBuffer := bytes.NewBuffer(blockBytes)
	syncMarker := [16]byte{'g', 'o', 'g', 'e', 'n', 'a', 'v', 'r', 'o', 'm', 'a', 'g', 'i', 'c', '1', '0'}

	// Write the header when we construct the writer
	header := &AvroContainerHeader{
		Magic: [4]byte{'O', 'b', 'j', 1},
		Meta: map[string][]byte{
			"avro.schema": []byte(DemoSchemaSchema),
			"avro.codec":  []byte(codec),
		},
		Sync: syncMarker,
	}

	err := header.Serialize(writer)
	if err != nil {
		return nil, err
	}

	avroWriter := &DemoSchemaContainerWriter{
		writer:          writer,
		syncMarker:      syncMarker,
		codec:           codec,
		recordsPerBlock: recordsPerBlock,
		blockBuffer:     blockBuffer,
	}

	if codec == Deflate {
		avroWriter.compressedWriter, err = flate.NewWriter(avroWriter.blockBuffer, flate.DefaultCompression)
		if err != nil {
			return nil, err
		}
	} else if codec == Snappy {
		avroWriter.compressedWriter = snappy.NewBufferedWriter(avroWriter.blockBuffer)
	} else if codec == Null {
		avroWriter.compressedWriter = avroWriter.blockBuffer
	}

	return avroWriter, nil
}

func (avroWriter *DemoSchemaContainerWriter) Flush() error {
	// Write out all of the buffered records as a new block
	// Must be called before closing to ensure the last block is written
	if fwWriter, ok := avroWriter.compressedWriter.(FlushableResettableWriter); ok {
		fwWriter.Flush()
		fwWriter.Reset(avroWriter.blockBuffer)
	}

	block := &AvroContainerBlock{
		NumRecords:  avroWriter.nextBlockRecords,
		RecordBytes: avroWriter.blockBuffer.Bytes(),
		Sync:        avroWriter.syncMarker,
	}
	err := block.Serialize(avroWriter.writer)
	if err != nil {
		return err
	}

	avroWriter.blockBuffer.Reset()
	avroWriter.nextBlockRecords = 0

	return nil
}

func (avroWriter *DemoSchemaContainerWriter) WriteRecord(record DemoSchema) error {
	// Serialize the new record into the compressed writer
	err := record.Serialize(avroWriter.compressedWriter)
	if err != nil {
		return err
	}
	avroWriter.nextBlockRecords += 1

	// If the block if full, flush and reset the compressed writer,
	// write the header and the block contents
	if avroWriter.nextBlockRecords >= avroWriter.recordsPerBlock {
		return avroWriter.Flush()
	}

	return nil
}
