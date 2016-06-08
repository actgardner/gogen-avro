package avro

import (
	"fmt"
	"io"
)

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
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

func readArrayInt(r io.Reader) ([]int32, error) {
	var err error
	var blkSize int64
	var arr []int32
	for {
		blkSize, err = readLong(r)
		if err != nil {
			return nil, err
		}
		if blkSize == 0 {
			break
		}
		if blkSize < 0 {
			blkSize = -blkSize
			_, err = readLong(r)
			if err != nil {
				return nil, err
			}
		}
		for i := int64(0); i < blkSize; i++ {
			elem, err := readInt(r)
			if err != nil {
				return nil, err
			}
			arr = append(arr, elem)
		}
	}
	return arr, nil
}

func readComplexUnionTestRecord(r io.Reader) (*ComplexUnionTestRecord, error) {
	var str ComplexUnionTestRecord
	var err error
	str.UnionField, err = readUnionNullArrayIntMapIntNestedUnionRecord(r)
	if err != nil {
		return nil, err
	}

	return &str, nil
}

func readInt(r io.Reader) (int32, error) {
	var v int
	buf := make([]byte, 1)
	for shift := uint(0); ; shift += 7 {
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}
		b := buf[0]
		v |= int(b&127) << shift
		if b&128 == 0 {
			break
		}
	}
	datum := (int32(v>>1) ^ -int32(v&1))
	return datum, nil
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

func readMapInt(r io.Reader) (map[string]int32, error) {
	m := make(map[string]int32)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		if blkSize == 0 {
			break
		}
		if blkSize < 0 {
			blkSize = -blkSize
			_, err := readLong(r)
			if err != nil {
				return nil, err
			}
		}
		for i := int64(0); i < blkSize; i++ {
			key, err := readString(r)
			if err != nil {
				return nil, err
			}
			val, err := readInt(r)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
	}
	return m, nil
}

func readNestedUnionRecord(r io.Reader) (*NestedUnionRecord, error) {
	var str NestedUnionRecord
	var err error
	str.IntField, err = readInt(r)
	if err != nil {
		return nil, err
	}

	return &str, nil
}

func readNull(_ io.Reader) (interface{}, error) {
	return nil, nil
}

func readString(r io.Reader) (string, error) {
	len, err := readLong(r)
	if err != nil {
		return "", err
	}
	bb := make([]byte, len)
	_, err = io.ReadFull(r, bb)
	if err != nil {
		return "", err
	}
	return string(bb), nil
}

func readUnionNullArrayIntMapIntNestedUnionRecord(r io.Reader) (UnionNullArrayIntMapIntNestedUnionRecord, error) {
	field, err := readLong(r)
	var unionStr UnionNullArrayIntMapIntNestedUnionRecord
	if err != nil {
		return unionStr, err
	}
	unionStr.UnionType = UnionNullArrayIntMapIntNestedUnionRecordTypeEnum(field)
	switch unionStr.UnionType {
	case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNull:
		val, err := readNull(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.Null = val
	case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumArrayInt:
		val, err := readArrayInt(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.ArrayInt = val
	case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumMapInt:
		val, err := readMapInt(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.MapInt = val
	case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNestedUnionRecord:
		val, err := readNestedUnionRecord(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.NestedUnionRecord = val

	default:
		return unionStr, fmt.Errorf("Invalid value for UnionNullArrayIntMapIntNestedUnionRecord")
	}
	return unionStr, nil
}

func writeArrayInt(r []int32, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for _, e := range r {
		err = writeInt(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeComplexUnionTestRecord(r *ComplexUnionTestRecord, w io.Writer) error {
	var err error
	err = writeUnionNullArrayIntMapIntNestedUnionRecord(r.UnionField, w)
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

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeMapInt(r map[string]int32, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeInt(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeNestedUnionRecord(r *NestedUnionRecord, w io.Writer) error {
	var err error
	err = writeInt(r.IntField, w)
	if err != nil {
		return err
	}

	return nil
}

func writeNull(_ interface{}, _ io.Writer) error {
	return nil
}

func writeString(r string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(r)
	} else {
		_, err = w.Write([]byte(r))
	}
	return err
}

func writeUnionNullArrayIntMapIntNestedUnionRecord(r UnionNullArrayIntMapIntNestedUnionRecord, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNull:
		return writeNull(r.Null, w)
	case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumArrayInt:
		return writeArrayInt(r.ArrayInt, w)
	case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumMapInt:
		return writeMapInt(r.MapInt, w)
	case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNestedUnionRecord:
		return writeNestedUnionRecord(r.NestedUnionRecord, w)

	}
	return fmt.Errorf("Invalid value for UnionNullArrayIntMapIntNestedUnionRecord")
}
