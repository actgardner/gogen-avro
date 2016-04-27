package generator

const writeIntMethod = `
func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := uint64((uint32(r) << 1) ^ uint32(r >> downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}
`

const encodeIntMethod = `
func encodeInt(w io.Writer, byteCount int, encoded uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil

}
`

type intField struct {
	name         string
	defaultValue int32
	hasDefault   bool
}

func (s *intField) Name() string {
	return toPublicName(s.name)
}

func (s *intField) FieldType() string {
	return "Int"
}

func (s *intField) GoType() string {
	return "int32"
}

func (s *intField) SerializerMethod() string {
	return "writeInt"
}

func (s *intField) AddStruct(p *Package) {}

func (s *intField) AddSerializer(p *Package) {
	p.addStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.addFunction(UTIL_FILE, "", "writeInt", writeIntMethod)
	p.addFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.addImport(UTIL_FILE, "io")
}
