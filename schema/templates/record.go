package templates

const RecordTemplate = `
{{ if ne .Doc "" }}
// {{ .Doc}}
{{ end }}  
type {{ .Name }} struct {
{{ range $i, $field := .Fields }}
	{{ if ne $field.Doc "" }}
	// {{ $field.Doc }}
	{{ end }}
	{{ if ne $field.Tags "" }}
		{{ $field.SimpleName }} {{ $field.Type.GoType }} ` + "`{{ $field.Tags }}`" + `
	{{ else }}
		{{ $field.SimpleName }} {{ $field.Type.GoType }}
	{{ end }}
{{ end }}
}

func {{ .ConstructorMethod }} ({{ .GoType}}) {
	return &{{ .Name }}{}
}

func Deserialize{{ .Name }}(r io.Reader) ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err	
	}
	return t, err
}

func Deserialize{{ .Name }}FromSchema(r io.Reader, schema string) ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err	
	}
	return t, err
}

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	var err error
	{{ range $i, $field := .Fields }}
	err = {{ .Type.SerializerMethod }}( r.{{ .GoName }}, w)
	if err != nil {
		return err			
	}
	{{ end }}
	return err
}

func (r {{ .GoType }}) Serialize(w io.Writer) error {
	return {{ .SerializerMethod }}(r, w)
}

func {{ .NewWriterMethod }}(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := {{ .ConstructorMethod }}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}

func (r {{ .GoType }}) Schema() string {
	return {{ printf "%q" .Schema }}
}

func (r {{ .GoType }}) SchemaName() string {
	return {{ printf "%q" .AvroName.String }}
}

func (_ {{ .GoType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetString(v string) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r {{ .GoType }}) Get(i int) types.Field {
	switch (i) {
	{{ range $i, $field := .Fields }}
	case {{ $i }}:
		{{ if $.ConstructableForField $field | ne "" }}
			{{ $.ConstructableForField $field }}	
		{{ end }}
		{{ if ne $field.Type.WrapperType "" }}
			return (*{{ $field.Type.WrapperType }})(&r.{{ $field.GoName }})
		{{ else }}
			return r.{{ $field.GoName }}
		{{ end }}
	{{ end }}
	}
	panic("Unknown field index")
}

func (r {{ .GoType }}) SetDefault(i int) {
	switch (i) {
	{{ range $i, $field := .Fields }}
        {{ if .HasDefault }}
	case {{ $i }}:
       	 	{{ $.DefaultForField $field }}
	{{ end }}
	{{ end }}
	}
	panic("Unknown field index")
}

func (_ {{ .GoType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) Finalize() { }

// container reader
type {{ .RecordReaderTypeName }} struct {
	r io.Reader
	p *vm.Program
}

func New{{ .RecordReaderTypeName }}(r io.Reader) (*{{ .RecordReaderTypeName }}, error){
	containerReader, err := container.NewReader(r)
	if err != nil {
		return nil, err
	}

	t := {{ .ConstructorMethod }}
	deser, err := compiler.CompileSchemaBytes([]byte(containerReader.AvroContainerSchema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	return &{{ .RecordReaderTypeName }} {
		r: containerReader,
		p: deser,
	}, nil
}

func (r {{ .RecordReaderTypeName }}) Read() ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}
        err := vm.Eval(r.r, r.p, t)
	return t, err
}
`
