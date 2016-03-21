package avro

import "io"
import "math"

type MapTestRecord struct {
	IntField    map[string]int32
	LongField   map[string]int64
	DoubleField map[string]float64
	StringField map[string]string
	FloatField  map[string]float32
	BoolField   map[string]bool
	BytesField  map[string][]byte
}

func (r MapTestRecord) Serialize(w io.Writer) error {
	return writeMapTestRecord(r, w)
}

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
	bb = make([]byte, 0, byteCount)
	for i := 0; i < byteCount; i++ {
		bb = append(bb, byte(bits&255))
		bits = bits >> 8
	}
	_, err = w.Write(bb)
	return err
}

func encodeInt(w io.Writer, byteCount int, encoded uint64) error {
	var err error
	var bb []byte
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	bb = make([]byte, 0, byteCount)

	if encoded == 0 {
		bb = append(bb, byte(0))
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			bb = append(bb, b)
		}
	}
	_, err = w.Write(bb)
	return err
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

func writeMapTestRecord(r MapTestRecord, w io.Writer) error {
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
