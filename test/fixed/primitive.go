package avro

import (
	"io"
)

func readFixedTestRecord(r io.Reader) (*FixedTestRecord, error) {
	var str FixedTestRecord
	var err error
	str.FixedField, err = readTestFixedType(r)
	if err != nil {
		return nil, err
	}

	return &str, nil
}

func readTestFixedType(r io.Reader) (TestFixedType, error) {
	var bb TestFixedType
	_, err := io.ReadFull(r, bb[:])
	return bb, err
}

func writeFixedTestRecord(r *FixedTestRecord, w io.Writer) error {
	var err error
	err = writeTestFixedType(r.FixedField, w)
	if err != nil {
		return err
	}

	return nil
}

func writeTestFixedType(r TestFixedType, w io.Writer) error {
	_, err := w.Write(r[:])
	return err
}
