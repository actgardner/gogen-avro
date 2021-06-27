// Package soe provides convenience methods to read and write Avro Single-Object Encoding headers
package soe

import (
	"encoding/binary"
	"errors"
	"io"
)

var HeaderV1 = []byte{0xC3, 0x01}

// WriteRecord writes the single-object encoding framing (the magic bytes and the Avro CRC64 fingerprint of the canonical form of the record schema), and the record itself into `w`.
func WriteRecord(w io.Writer, record AvroRecord) error {
	fp := HeaderV1
	if err := binary.Write(w, binary.LittleEndian, fp); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, record.AvroCRC64Fingerprint()); err != nil {
		return err
	}

	return record.Serialize(w)
}

// ReadHeader reads the magic bytes and CRC64 schema fingerprint from the given reader and returns the fingerprint.
func ReadHeader(r io.Reader) ([]byte, error) {
	b := make([]byte, 2)
	n, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	if b[0] != HeaderV1[0] || b[1] != HeaderV1[1] {
		return nil, errors.New("avro v1 header invalid")
	}

	var header []byte

	header = make([]byte, 8)

	n, err = r.Read(header)

	if n != 8 || err != nil {
		return nil, errors.New("fingerprint header not read")
	}

	return header, nil
}
