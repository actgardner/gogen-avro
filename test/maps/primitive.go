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

func readMapBool(r io.Reader) (map[string]bool, error) {
	m := make(map[string]bool)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Decoding block size \n", blkSize)
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
			val, err := readBool(r)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
	}
	return m, nil
}

func readMapBytes(r io.Reader) (map[string][]byte, error) {
	m := make(map[string][]byte)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Decoding block size \n", blkSize)
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
			val, err := readBytes(r)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
	}
	return m, nil
}

func readMapDouble(r io.Reader) (map[string]float64, error) {
	m := make(map[string]float64)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Decoding block size \n", blkSize)
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
			val, err := readDouble(r)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
	}
	return m, nil
}

func readMapFloat(r io.Reader) (map[string]float32, error) {
	m := make(map[string]float32)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Decoding block size \n", blkSize)
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
			val, err := readFloat(r)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
	}
	return m, nil
}

func readMapInt(r io.Reader) (map[string]int32, error) {
	m := make(map[string]int32)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Decoding block size \n", blkSize)
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

func readMapLong(r io.Reader) (map[string]int64, error) {
	m := make(map[string]int64)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Decoding block size \n", blkSize)
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
			val, err := readLong(r)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
	}
	return m, nil
}

func readMapString(r io.Reader) (map[string]string, error) {
	m := make(map[string]string)
	for {
		blkSize, err := readLong(r)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Decoding block size \n", blkSize)
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
			val, err := readString(r)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
	}
	return m, nil
}

func readMapTestRecord(r io.Reader) (*MapTestRecord, error) {
	var str MapTestRecord
	var err error
	str.IntField, err = readMapInt(r)
	if err != nil {
		return nil, err
	}
	str.LongField, err = readMapLong(r)
	if err != nil {
		return nil, err
	}
	str.DoubleField, err = readMapDouble(r)
	if err != nil {
		return nil, err
	}
	str.StringField, err = readMapString(r)
	if err != nil {
		return nil, err
	}
	str.FloatField, err = readMapFloat(r)
	if err != nil {
		return nil, err
	}
	str.BoolField, err = readMapBool(r)
	if err != nil {
		return nil, err
	}
	str.BytesField, err = readMapBytes(r)
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

func writeMapBool(r map[string]bool, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeBool(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapBytes(r map[string][]byte, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeBytes(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapDouble(r map[string]float64, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeDouble(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapFloat(r map[string]float32, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeFloat(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
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

func writeMapLong(r map[string]int64, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeLong(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapString(r map[string]string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeString(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapTestRecord(r *MapTestRecord, w io.Writer) error {
	var err error
	err = writeMapInt(r.IntField, w)
	if err != nil {
		return err
	}
	err = writeMapLong(r.LongField, w)
	if err != nil {
		return err
	}
	err = writeMapDouble(r.DoubleField, w)
	if err != nil {
		return err
	}
	err = writeMapString(r.StringField, w)
	if err != nil {
		return err
	}
	err = writeMapFloat(r.FloatField, w)
	if err != nil {
		return err
	}
	err = writeMapBool(r.BoolField, w)
	if err != nil {
		return err
	}
	err = writeMapBytes(r.BytesField, w)
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
