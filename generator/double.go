package generator

const writeDoubleMethod = `
func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}
`

const readDoubleMethod = `
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
`

type doubleField struct {
	name         string
	defaultValue float64
	hasDefault   bool
}

func (s *doubleField) Name() string {
	return toPublicName(s.name)
}

func (s *doubleField) FieldType() string {
	return "Double"
}

func (s *doubleField) GoType() string {
	return "float64"
}

func (s *doubleField) SerializerMethod() string {
	return "writeDouble"
}

func (s *doubleField) DeserializerMethod() string {
	return "readDouble"
}

func (s *doubleField) AddStruct(*Package) {}

func (s *doubleField) AddSerializer(p *Package) {
	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addFunction(UTIL_FILE, "", "writeDouble", writeDoubleMethod)
	p.addFunction(UTIL_FILE, "", "encodeFloat", encodeFloatMethod)
	p.addImport(UTIL_FILE, "io")
	p.addImport(UTIL_FILE, "math")
}

func (s *doubleField) AddDeserializer(p *Package) {
	p.addFunction(UTIL_FILE, "", "readDouble", readDoubleMethod)
	p.addImport(UTIL_FILE, "io")
	p.addImport(UTIL_FILE, "math")
	p.addImport(UTIL_FILE, "encoding/binary")
}
