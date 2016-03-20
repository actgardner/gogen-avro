package generator

const byteWriterInterface = `
type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
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

type boolField struct {
	name         string
	defaultValue bool
	hasDefault   bool
}

func (s *boolField) Name() string {
	return toPublicName(s.name)
}

func (s *boolField) FieldType() string {
	return "Bool"
}

func (s *boolField) GoType() string {
	return "bool"
}

func (s *boolField) AuxStructs(aux map[string]string, _ map[string]string) {
	aux["ByteWriter"] = byteWriterInterface
	aux["writeBool"] = writeBoolMethod
}

func (s *boolField) SerializerMethod() string {
	return "writeBool"
}
