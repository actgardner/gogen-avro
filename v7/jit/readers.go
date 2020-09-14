package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/actgardner/gogen-avro/v7/vm/types"
)

type ByteReader interface {
	ReadByte() (byte, error)
}

func assignBoolToBool(r io.Reader, t *bool) error {
	v, err := readBool(r)
	if err != nil {
		return err
	}

	*t = v
	return nil
}

func assignIntToInt(r io.Reader, t *int32) error {
	v, err := readInt(r)
	if err != nil {
		return err
	}
	*t = v
	return nil
}

func assignLongToLong(r io.Reader, t *int64) error {
	v, err := readLong(r)
	if err != nil {
		return err
	}
	*t = v
	return nil
}

func assignFloatToFloat(r io.Reader, t *float32) error {
	v, err := readFloat(r)
	if err != nil {
		return err
	}
	*t = v
	return nil
}

func assignDoubleToDouble(r io.Reader, t *float64) error {
	v, err := readDouble(r)
	if err != nil {
		return err
	}
	*t = v
	return nil
}

func assignBytesToBytes(r io.Reader, t *[]byte) error {
	v, err := readBytes(r)
	if err != nil {
		return err
	}
	*t = v
	return nil
}

func assignStringToString(r io.Reader, t *string) error {
	v, err := readString(r)
	if err != nil {
		return err
	}
	*t = v
	return nil
}

func readBool(r io.Reader) (bool, error) {
	var b byte
	var err error
	if br, ok := r.(ByteReader); ok {
		b, err = br.ReadByte()
	} else {
		bs := make([]byte, 1)
		_, err = io.ReadFull(r, bs)
		if err != nil {
			return false, err
		}
		b = bs[0]
	}
	return b == 1, nil
}

func readBytes(r io.Reader) ([]byte, error) {
	size, err := readLong(r)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return []byte{}, nil
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
	var b byte
	var err error
	if br, ok := r.(ByteReader); ok {
		for shift := uint(0); ; shift += 7 {
			if b, err = br.ReadByte(); err != nil {
				return 0, err
			}
			v |= int(b&127) << shift
			if b&128 == 0 {
				break
			}
		}
	} else {
		buf := make([]byte, 1)
		for shift := uint(0); ; shift += 7 {
			if _, err := io.ReadFull(r, buf); err != nil {
				return 0, err
			}
			b = buf[0]
			v |= int(b&127) << shift
			if b&128 == 0 {
				break
			}
		}
	}
	datum := (int32(v>>1) ^ -int32(v&1))
	return datum, nil
}

func readLong(r io.Reader) (int64, error) {
	var v uint64
	var b byte
	var err error
	if br, ok := r.(ByteReader); ok {
		for shift := uint(0); ; shift += 7 {
			if b, err = br.ReadByte(); err != nil {
				return 0, err
			}
			v |= uint64(b&127) << shift
			if b&128 == 0 {
				break
			}
		}
	} else {
		buf := make([]byte, 1)
		for shift := uint(0); ; shift += 7 {
			if _, err = io.ReadFull(r, buf); err != nil {
				return 0, err
			}
			b = buf[0]
			v |= uint64(b&127) << shift
			if b&128 == 0 {
				break
			}
		}
	}
	datum := (int64(v>>1) ^ -int64(v&1))
	return datum, nil
}

func readString(r io.Reader) (string, error) {
	len, err := readLong(r)
	if err != nil {
		return "", err
	}

	// makeslice can fail depending on available memory.
	// We arbitrarily limit string size to sane default (~2.2GB).
	if len < 0 || len > math.MaxInt32 {
		return "", fmt.Errorf("string length out of range: %d", len)
	}

	if len == 0 {
		return "", nil
	}

	bb := make([]byte, len)
	_, err = io.ReadFull(r, bb)
	if err != nil {
		return "", err
	}
	return string(bb), nil
}

func readFixed(r io.Reader, size int) ([]byte, error) {
	bb := make([]byte, size)
	_, err := io.ReadFull(r, bb)
	return bb, err
}
