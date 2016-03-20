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

func (s *doubleField) AuxStructs(aux map[string]string, imports map[string]string) {
	imports["math"] = "import \"math\""
	aux["writeDouble"] = writeDoubleMethod
	aux["encodeFloat"] = encodeFloatMethod
	aux["ByteWriter"] = byteWriterInterface
}

func (s *doubleField) SerializerMethod() string {
	return "writeDouble"
}
