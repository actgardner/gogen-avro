package avro

import (
	"io"
)

type AvroContainerBlock struct {
	NumRecords  int64
	RecordBytes []byte
	Sync        Sync
}

func DeserializeAvroContainerBlock(r io.Reader) (*AvroContainerBlock, error) {
	return readAvroContainerBlock(r)
}

func (r AvroContainerBlock) Serialize(w io.Writer) error {
	return writeAvroContainerBlock(&r, w)
}
