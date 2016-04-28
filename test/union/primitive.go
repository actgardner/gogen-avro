package avro

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
}

func encodeFloat(w io.Writer, byteCount int, bits uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}
	for i := 0; i < byteCount; i++ {
		if bw != nil {
			err = bw.WriteByte(byte(bits & 255))
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(bits&255))
		}
		bits = bits >> 8
	}
	if bw == nil {
		_, err = w.Write(bb)
		return err
	}
	return nil
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

func readBool(r io.Reader) (bool, error) {
	b := make([]byte, 1)
	_, err := r.Read(b)
	if err != nil {
		return false, err
	}
	return b[0] == 1, nil
}

func readBytes(r io.Reader) ([]byte, error) {
	size, err := readLong(r)
	if err != nil {
		return nil, err
	}
	bb := make([]byte, size)
	_, err = io.ReadFull(r, bb)
	return bb, err
}

func readDouble(r io.Reader) (float64, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint64(buf)
	val := math.Float64frombits(bits)
	return val, nil
}

func readFloat(r io.Reader) (float32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint32(buf)
	val := math.Float32frombits(bits)
	return val, nil

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

func readNull(_ io.Reader) (interface{}, error) {
	return nil, nil
}

func readPrimitiveUnionTestRecord(r io.Reader) (*PrimitiveUnionTestRecord, error) {
	var str PrimitiveUnionTestRecord
	var err error
	str.UnionField, err = readUnionIntLongFloatDoubleStringBoolBytesNull(r)
	if err != nil {
		return nil, err
	}

	return &str, nil
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

func readUnionIntLongFloatDoubleStringBoolBytesNull(r io.Reader) (UnionIntLongFloatDoubleStringBoolBytesNull, error) {
	field, err := readLong(r)
	var unionStr UnionIntLongFloatDoubleStringBoolBytesNull
	if err != nil {
		return unionStr, err
	}
	unionStr.UnionType = UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum(field)
	switch unionStr.UnionType {
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumInt:
		val, err := readInt(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.Int = val
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumLong:
		val, err := readLong(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.Long = val
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumFloat:
		val, err := readFloat(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.Float = val
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumDouble:
		val, err := readDouble(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.Double = val
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumString:
		val, err := readString(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.String = val
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumBool:
		val, err := readBool(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.Bool = val
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumBytes:
		val, err := readBytes(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.Bytes = val
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumNull:
		val, err := readNull(r)
		if err != nil {
			return unionStr, err
		}
		unionStr.Null = val

	default:
		return unionStr, fmt.Errorf("Invalid value for UnionIntLongFloatDoubleStringBoolBytesNull")
	}
	return unionStr, nil
}

func writeBool(r bool, w io.Writer) error {
	var b byte
	if r {
		b = byte(1)
	}

	var err error
	if bw, ok := w.(ByteWriter); ok {
		err = bw.WriteByte(b)
	} else {
		bb := make([]byte, 1)
		bb[0] = b
		_, err = w.Write(bb)
	}
	if err != nil {
		return err
	}
	return nil
}

func writeBytes(r []byte, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	_, err = w.Write(r)
	return err
}

func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}

func writeFloat(r float32, w io.Writer) error {
	bits := uint64(math.Float32bits(r))
	const byteCount = 4
	return encodeFloat(w, byteCount, bits)
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

func writeNull(_ interface{}, _ io.Writer) error {
	return nil
}

func writePrimitiveUnionTestRecord(r *PrimitiveUnionTestRecord, w io.Writer) error {
	var err error
	err = writeUnionIntLongFloatDoubleStringBoolBytesNull(r.UnionField, w)
	if err != nil {
		return err
	}

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

func writeUnionIntLongFloatDoubleStringBoolBytesNull(r UnionIntLongFloatDoubleStringBoolBytesNull, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumInt:
		return writeInt(r.Int, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumLong:
		return writeLong(r.Long, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumFloat:
		return writeFloat(r.Float, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumDouble:
		return writeDouble(r.Double, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumString:
		return writeString(r.String, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumBool:
		return writeBool(r.Bool, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumBytes:
		return writeBytes(r.Bytes, w)
	case UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumNull:
		return writeNull(r.Null, w)

	}
	return fmt.Errorf("Invalid value for UnionIntLongFloatDoubleStringBoolBytesNull")
}
