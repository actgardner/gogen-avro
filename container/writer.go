package container

import (
	"bytes"
	"compress/flate"
	"io"
)

type Codec string

const (
	Null    Codec = "null"
	Deflate Codec = "deflate"
	Snappy  Codec = "snappy"
)

type CloseableResettableWriter interface {
	Close() error
	Reset(io.Writer)
}

type Writer struct {
	writer           io.Writer
	syncMarker       [16]byte
	codec            Codec
	recordsPerBlock  int64
	blockBuffer      *bytes.Buffer
	compressedWriter io.Writer
	nextBlockRecords int64
	headerWritten    bool
}

func NewWriter(writer io.Writer, codec Codec, recordsPerBlock int64) (*Writer, error) {
	blockBytes := make([]byte, 0)
	blockBuffer := bytes.NewBuffer(blockBytes)

	avroWriter := &Writer{
		writer:          writer,
		syncMarker:      [16]byte{'g', 'o', 'g', 'e', 'n', 'a', 'v', 'r', 'o', 'm', 'a', 'g', 'i', 'c', '1', '0'},
		codec:           codec,
		recordsPerBlock: recordsPerBlock,
		blockBuffer:     blockBuffer,
		headerWritten:   false,
	}
	var err error
	if codec == Deflate {
		avroWriter.compressedWriter, err = flate.NewWriter(avroWriter.blockBuffer, flate.DefaultCompression)
		if err != nil {
			return nil, err
		}
	} else if codec == Snappy {
		avroWriter.compressedWriter = newSnappyWriter(avroWriter.blockBuffer)
	} else if codec == Null {
		avroWriter.compressedWriter = avroWriter.blockBuffer
	}

	return avroWriter, nil
}

func (avroWriter *Writer) writeHeader(schema string) error {
	header := &AvroContainerHeader{
		Magic: [4]byte{'O', 'b', 'j', 1},
		Meta: map[string][]byte{
			"avro.schema": []byte(schema),
			"avro.codec":  []byte(avroWriter.codec),
		},
		Sync: avroWriter.syncMarker,
	}
	return header.Serialize(avroWriter.writer)
}

func (avroWriter *Writer) WriteRecord(record AvroRecord) error {
	var err error
	// Lazily write the header when the first record is written
	if !avroWriter.headerWritten {
		avroWriter.headerWritten = true
		err = avroWriter.writeHeader(record.Schema())
		if err != nil {
			return err
		}
	}
	// Serialize the new record into the compressed writer
	err = record.Serialize(avroWriter.compressedWriter)
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

func (avroWriter *Writer) Flush() error {
	// Write out all of the buffered records as a new block
	// Must be called before closing to ensure the last block is written
	if fwWriter, ok := avroWriter.compressedWriter.(CloseableResettableWriter); ok {
		fwWriter.Close()
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
