package types

import (
	"fmt"
	"github.com/actgardner/gogen-avro/generator"
)

const writeFloatMethod = `
func writeFloat(r float32, w io.Writer) error {
	bits := uint64(math.Float32bits(r))
	const byteCount = 4
	return encodeFloat(w, byteCount, bits)
}
`

const encodeFloatMethod = `
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
`
const readFloatMethod = `
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
`

type floatField struct {
	primitiveField
}

func NewFloatField(definition interface{}) *floatField {
	return &floatField{primitiveField{
		definition:         definition,
		name:               "Float",
		goType:             "float32",
		serializerMethod:   "writeFloat",
		deserializerMethod: "readFloat",
	}}
}

func (e *floatField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeFloat", writeFloatMethod)
	p.AddFunction(UTIL_FILE, "", "encodeFloat", encodeFloatMethod)
	p.AddImport(UTIL_FILE, "math")
	p.AddImport(UTIL_FILE, "io")
}

func (e *floatField) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readFloat", readFloatMethod)
	p.AddImport(UTIL_FILE, "math")
	p.AddImport(UTIL_FILE, "encoding/binary")
	p.AddImport(UTIL_FILE, "io")
}

func (s *floatField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected float as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}
