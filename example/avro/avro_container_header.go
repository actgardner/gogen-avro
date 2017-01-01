package avro

import (
	"io"
)

type AvroContainerHeader struct {
	Magic Magic
	Meta  map[string][]byte
	Sync  Sync
}

func DeserializeAvroContainerHeader(r io.Reader) (*AvroContainerHeader, error) {
	return readAvroContainerHeader(r)
}

func (r AvroContainerHeader) Serialize(w io.Writer) error {
	return writeAvroContainerHeader(&r, w)
}
