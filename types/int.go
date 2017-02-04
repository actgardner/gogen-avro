package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

const writeIntMethod = `
func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := uint64((uint32(r) << 1) ^ uint32(r >> downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}
`

const encodeIntMethod = `
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
`

const readIntMethod = `
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
`

type intField struct {
	name         string
	defaultValue int32
	hasDefault   bool
}

func (s *intField) HasDefault() bool {
	return s.hasDefault
}

func (s *intField) Default() interface{} {
	return s.defaultValue
}

func (s *intField) AvroName() string {
	return s.name
}

func (s *intField) GoName() string {
	return generator.ToPublicName(s.name)
}

func (s *intField) FieldType() string {
	return "Int"
}

func (s *intField) GoType() string {
	return "int32"
}

func (s *intField) SerializerMethod() string {
	return "writeInt"
}

func (s *intField) DeserializerMethod() string {
	return "readInt"
}

func (s *intField) AddStruct(p *generator.Package) {}

func (s *intField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeInt", writeIntMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *intField) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readInt", readIntMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *intField) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *intField) Schema(names map[QualifiedName]interface{}) interface{} {
	return "int"
}
