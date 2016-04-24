package avro

import "fmt"
import "io"

type ComplexUnionTestRecord struct {
	UnionField UnionNullArrayIntMapIntNestedUnionRecord
}

func (r ComplexUnionTestRecord) Serialize(w io.Writer) error {
	return writeComplexUnionTestRecord(&r, w)
}

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type NestedUnionRecord struct {
	IntField int32
}

func (r NestedUnionRecord) Serialize(w io.Writer) error {
	return writeNestedUnionRecord(&r, w)
}

type StringWriter interface {
	WriteString(string) (int, error)
}

type UnionNullArrayIntMapIntNestedUnionRecord struct {
	Null              interface{}
	ArrayInt          []int32
	MapInt            map[string]int32
	NestedUnionRecord *NestedUnionRecord
	UnionType         UnionNullArrayIntMapIntNestedUnionRecordTypeEnum
}

type UnionNullArrayIntMapIntNestedUnionRecordTypeEnum int

const (
	UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNull              UnionNullArrayIntMapIntNestedUnionRecordTypeEnum = 0
	UnionNullArrayIntMapIntNestedUnionRecordTypeEnumArrayInt          UnionNullArrayIntMapIntNestedUnionRecordTypeEnum = 1
	UnionNullArrayIntMapIntNestedUnionRecordTypeEnumMapInt            UnionNullArrayIntMapIntNestedUnionRecordTypeEnum = 2
	UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNestedUnionRecord UnionNullArrayIntMapIntNestedUnionRecordTypeEnum = 3
)

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
