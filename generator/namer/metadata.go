package namer

var MetadataKey = struct{}{}

type TypeMetadata struct {
	Name              string
	GoType            string
	SerializerMethod  string
	ConstructorMethod string
	WrapperType       string
}

type RecordMetadata struct {
	TypeMetadata

	SerializerMethod     string
	NewWriterMethod      string
	ConstructorMethod    string
	RecordReaderTypeName string
}

type FixedMetadata struct {
	TypeMetadata

	SerializerMethod string
}

type EnumMetadata struct {
	TypeMetadata

	SerializerMethod string
	FromStringMethod string
	SymbolNames      []string
}

type FieldMetadata struct {
	Name string
}
