package generator

const writeLongMethod = `
func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}
`

const readLongMethod = `
func readLong(r io.Reader) (int64, error) {
	var v uint64
	buf := make([]byte, 1)
	for shift := uint(0); ; shift += 7 {
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}
		b := buf[0]
		v |= uint64(b&127) << shift
		if b&128 == 0 {
			break
		}
	}
	datum := (int64(v>>1) ^ -int64(v&1))
	return datum, nil
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

func (s *longField) SerializerMethod() string {
	return "writeLong"
}

func (s *longField) DeserializerMethod() string {
	return "readLong"
}

func (s *longField) AddStruct(p *Package) {}

func (s *longField) AddSerializer(p *Package) {
	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.addFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.addImport(UTIL_FILE, "io")
}

func (s *longField) AddDeserializer(p *Package) {
	p.addFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.addImport(UTIL_FILE, "io")
}
