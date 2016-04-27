package generator

const writeBytesMethod = `
func writeBytes(r []byte, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	_, err = w.Write(r)
	return err
}
`

type bytesField struct {
	name         string
	defaultValue []byte
	hasDefault   bool
}

func (s *bytesField) Name() string {
	return toPublicName(s.name)
}

func (s *bytesField) FieldType() string {
	return "Bytes"
}

func (s *bytesField) GoType() string {
	return "[]byte"
}

func (s *bytesField) SerializerMethod() string {
	return "writeBytes"
}

func (s *bytesField) AddStruct(*Package) {}

func (s *bytesField) AddSerializer(p *Package) {
	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addFunction(UTIL_FILE, "", "writeBytes", writeBytesMethod)
	p.addFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.addFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.addImport(UTIL_FILE, "io")
}
