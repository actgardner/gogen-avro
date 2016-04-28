package avro

import (
	"fmt"
	"io"
)

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
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

func readLong(r io.Reader) (int64, error) {
	var v uint64
	buf := make([]byte, 1)
	for shift := uint(0); ; shift += 7 {
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}
		b := buf[0]
		v |= uint64(b&127) << shift
		if b&128 == 0 {
			break
		}
	}
	datum := (int64(v>>1) ^ -int64(v&1))
	return datum, nil
}

func readNull(_ io.Reader) (interface{}, error) {
	return nil, nil
}

func readRecursiveUnionTestRecord(r io.Reader) (*RecursiveUnionTestRecord, error) {
	var str RecursiveUnionTestRecord
	var err error
	str.RecursiveField, err = readUnionNullRecursiveUnionTestRecord(r)
	if err != nil {
		return nil, err
	}

	return &str, nil
}

func readUnionNullRecursiveUnionTestRecord(r io.Reader) (UnionNullRecursiveUnionTestRecord, error) {
	field, err := readLong(r)
	var unionStr UnionNullRecursiveUnionTestRecord
	if err != nil {
		return unionStr, err
	}
	unionStr.UnionType = UnionNullRecursiveUnionTestRecordTypeEnum(field)
	switch unionStr.UnionType {
	case UnionNullRecursiveUnionTestRecordTypeEnumNull:
		val, err := readNull(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.Null = val
	case UnionNullRecursiveUnionTestRecordTypeEnumRecursiveUnionTestRecord:
		val, err := readRecursiveUnionTestRecord(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.RecursiveUnionTestRecord = val

	default:
		return unionStr, fmt.Errorf("Invalid value for UnionNullRecursiveUnionTestRecord")
	}
	return unionStr, nil
}

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeNull(_ interface{}, _ io.Writer) error {
	return nil
}

func writeRecursiveUnionTestRecord(r *RecursiveUnionTestRecord, w io.Writer) error {
	var err error
	err = writeUnionNullRecursiveUnionTestRecord(r.RecursiveField, w)
	if err != nil {
		return err
	}

	return nil
}

func writeUnionNullRecursiveUnionTestRecord(r UnionNullRecursiveUnionTestRecord, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullRecursiveUnionTestRecordTypeEnumNull:
		return writeNull(r.Null, w)
	case UnionNullRecursiveUnionTestRecordTypeEnumRecursiveUnionTestRecord:
		return writeRecursiveUnionTestRecord(r.RecursiveUnionTestRecord, w)

	}
	return fmt.Errorf("Invalid value for UnionNullRecursiveUnionTestRecord")
}
