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

func (s *bytesField) AuxStructs(aux map[string]string, imports map[string]string) {
	aux["writeBytes"] = writeBytesMethod
	aux["writeLong"] = writeLongMethod
	aux["encodeInt"] = encodeIntMethod
	aux["ByteWriter"] = byteWriterInterface

}

func (s *bytesField) SerializerMethod() string {
	return "writeBytes"
}
