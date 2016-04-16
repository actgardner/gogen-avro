package avro

import "io"

type EnumTestRecord struct {
	EnumField TestEnumTypeEnum
}

func (r EnumTestRecord) Serialize(w io.Writer) error {
	return writeEnumTestRecord(&r, w)
}

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type TestEnumTypeEnum int32

const (
	TestSymbol1 TestEnumTypeEnum = 0
	TestSymbol2 TestEnumTypeEnum = 1
	TestSymbol3 TestEnumTypeEnum = 2
)

func (e TestEnumTypeEnum) String() string {
	switch e {
	case TestSymbol1:
		return "TestSymbol1"
	case TestSymbol2:
		return "testSymbol2"
	case TestSymbol3:
		return "testSymbol3"

	}
	return "Unknown"
}

func encodeInt(w io.Writer, byteCount int, encoded uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil

}

func writeEnumTestRecord(r *EnumTestRecord, w io.Writer) error {
	var err error
	err = writeTestEnumTypeEnum(r.EnumField, w)
	if err != nil {
		return err
	}

	return nil
}

func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := uint64((uint32(r) << 1) ^ uint32(r>>downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}

func writeTestEnumTypeEnum(r TestEnumTypeEnum, w io.Writer) error {
	return writeInt(int32(r), w)
}
