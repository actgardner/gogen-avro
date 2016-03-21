package generator

const writeFloatMethod = `
func writeFloat(r float32, w io.Writer) error {
	bits := uint64(math.Float32bits(r))
	const byteCount = 4
	return encodeFloat(w, byteCount, bits)
}
`

const encodeFloatMethod = `
func encodeFloat(w io.Writer, byteCount int, bits uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}
	for i := 0; i < byteCount; i++ {
		if bw != nil {
			err = bw.WriteByte(byte(bits & 255))
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(bits&255))
		}
		bits = bits >> 8
	}
	if bw == nil {
		_, err = w.Write(bb)
		return err
	}
	return nil
}
`

type floatField struct {
	name         string
	defaultValue float32
	hasDefault   bool
}

func (s *floatField) Name() string {
	return toPublicName(s.name)
}

func (s *floatField) FieldType() string {
	return "Float"
}

func (s *floatField) GoType() string {
	return "float32"
}

func (s *floatField) SerializerNs(imports, aux map[string]string) {
	imports["math"] = "import \"math\""
	aux["writeFloat"] = writeFloatMethod
	aux["encodeFloat"] = encodeFloatMethod
	aux["ByteWriter"] = byteWriterInterface
}

func (s *floatField) SerializerMethod() string {
	return "writeFloat"
}
