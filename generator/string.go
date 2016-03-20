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

func (s *stringField) AuxStructs(aux map[string]string, _ map[string]string) {
	aux["writeLong"] = writeLongMethod
	aux["writeString"] = writeStringMethod
	aux["encodeInt"] = encodeIntMethod
	aux["ByteWriter"] = byteWriterInterface
	aux["StringWriter"] = stringWriterInterface
}

func (s *stringField) SerializerMethod() string {
	return "writeString"
}
