package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

const byteWriterInterface = `
type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
} 
`

const byteReaderInterface = `
type ByteReader interface {
	ReadByte() (byte, error)
} 
`

const writeBoolMethod = `

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
`

const readBoolMethod = `
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
`

type boolField struct {
	name         string
	defaultValue bool
	hasDefault   bool
}

func (s *boolField) Name() string {
	return generator.ToPublicName(s.name)
}

func (s *boolField) FieldType() string {
	return "Bool"
}

func (s *boolField) GoType() string {
	return "bool"
}

func (s *boolField) SerializerMethod() string {
	return "writeBool"
}

func (s *boolField) DeserializerMethod() string {
	return "readBool"
}

func (s *boolField) AddStruct(*generator.Package) {}

func (s *boolField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeBool", writeBoolMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *boolField) AddDeserializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteReader", byteReaderInterface)
	p.AddFunction(UTIL_FILE, "", "readBool", readBoolMethod)
	p.AddImport(UTIL_FILE, "io")
}
