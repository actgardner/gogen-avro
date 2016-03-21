package generator

const writeLongMethod = `
func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}
`

type longField struct {
	name         string
	defaultValue int64
	hasDefault   bool
}

func (s *longField) Name() string {
	return toPublicName(s.name)
}

func (s *longField) FieldType() string {
	return "Long"
}

func (s *longField) GoType() string {
	return "int64"
}

func (s *longField) SerializerNs(imports, aux map[string]string) {
	aux["writeLong"] = writeLongMethod
	aux["encodeInt"] = encodeIntMethod
	aux["ByteWriter"] = byteWriterInterface
}

func (s *longField) SerializerMethod() string {
	return "writeLong"
}
