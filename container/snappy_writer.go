package container

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/snappy"
	"hash/crc32"
	"io"
)

// A Writer that buffers until it's closed, then
// emits one Snappy-encoded block with the CRC suffix
// required by the Avro spec
type snappyWriter struct {
	writer      io.Writer
	inputBuffer *bytes.Buffer
	outputBytes []byte
}

func newSnappyWriter(writer io.Writer) *snappyWriter {
	return &snappyWriter{
		writer:      writer,
		inputBuffer: bytes.NewBuffer(make([]byte, 0)),
		outputBytes: make([]byte, 0),
	}
}

func (w *snappyWriter) Write(buf []byte) (int, error) {
	return w.inputBuffer.Write(buf)
}

func (w *snappyWriter) Close() error {
	w.outputBytes = snappy.Encode(w.outputBytes, w.inputBuffer.Bytes())
	_, err := w.writer.Write(w.outputBytes)
	if err != nil {
		return err
	}
	return binary.Write(w.writer, binary.BigEndian, crc32.ChecksumIEEE(w.inputBuffer.Bytes()))
}

func (w *snappyWriter) Reset(writer io.Writer) {
	w.outputBytes = w.outputBytes[:0]
	w.inputBuffer.Reset()
	w.writer = writer
}
