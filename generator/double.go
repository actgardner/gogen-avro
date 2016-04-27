package generator

const writeDoubleMethod = `
func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
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

func (s *doubleField) AddStruct(*Package) {}

func (s *doubleField) AddSerializer(p *Package) {
	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addFunction(UTIL_FILE, "", "writeDouble", writeDoubleMethod)
	p.addFunction(UTIL_FILE, "", "encodeFloat", encodeFloatMethod)
	p.addImport(UTIL_FILE, "io")
}
