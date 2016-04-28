package avro

import (
	"io"
)

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
