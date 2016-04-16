package avro

import "io"

type FixedTestRecord struct {
	FixedField TestFixedType
}

func (r FixedTestRecord) Serialize(w io.Writer) error {
	return writeFixedTestRecord(&r, w)
}

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type TestFixedType [12]byte

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
