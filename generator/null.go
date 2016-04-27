package generator

const writeNullMethod = `
func writeNull(_ interface{}, _ io.Writer) error {
	return nil
}
`

type nullField struct {
	name       string
	hasDefault bool
}

func (s *nullField) Name() string {
	return toPublicName(s.name)
}

func (s *nullField) FieldType() string {
	return "Null"
}

func (s *nullField) GoType() string {
	return "interface{}"
}

func (s *nullField) SerializerNs(imports, aux map[string]string) {
	aux["writeNull"] = writeNullMethod
}

func (s *nullField) SerializerMethod() string {
	return "writeNull"
}

func (s *nullField) AddStruct(p *Package) {}

func (s *nullField) AddSerializer(p *Package) {
	p.addFunction(UTIL_FILE, "", "writeNull", writeNullMethod)
	p.addImport(UTIL_FILE, "io")
}
