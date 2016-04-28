package generator

const stringWriterInterface = `
type StringWriter interface {
	WriteString(string) (int, error)
}
`

const writeStringMethod = `
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
`

const readStringMethod = `
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
`

type stringField struct {
	name         string
	defaultValue string
	hasDefault   bool
}

func (s *stringField) Name() string {
	return toPublicName(s.name)
}

func (s *stringField) FieldType() string {
	return "String"
}

func (s *stringField) GoType() string {
	return "string"
}

func (s *stringField) SerializerMethod() string {
	return "writeString"
}

func (s *stringField) DeserializerMethod() string {
	return "readString"
}

func (s *stringField) AddStruct(*Package) {}

func (s *stringField) AddSerializer(p *Package) {
	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addStruct(UTIL_FILE, "StringWriter", stringWriterInterface)
	p.addFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.addFunction(UTIL_FILE, "", "writeString", writeStringMethod)
	p.addFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.addImport(UTIL_FILE, "io")
}

func (s *stringField) AddDeserializer(p *Package) {
	p.addFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.addFunction(UTIL_FILE, "", "readString", readStringMethod)
	p.addImport(UTIL_FILE, "io")
}
