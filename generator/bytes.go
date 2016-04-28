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

const readBytesMethod = `
func readBytes(r io.Reader) ([]byte, error) {
	size, err := readLong(r)
	if err != nil {
		return nil, err
	}
	bb := make([]byte, size)
	_, err = io.ReadFull(r, bb)
	return bb, err
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

func (s *bytesField) DeserializerMethod() string {
	return "readBytes"
}

func (s *bytesField) AddStruct(*Package) {}

func (s *bytesField) AddSerializer(p *Package) {
	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addFunction(UTIL_FILE, "", "writeBytes", writeBytesMethod)
	p.addFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.addFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.addImport(UTIL_FILE, "io")
}

func (s *bytesField) AddDeserializer(p *Package) {
	p.addFunction(UTIL_FILE, "", "readBytes", readBytesMethod)
	p.addFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.addImport(UTIL_FILE, "io")
}
