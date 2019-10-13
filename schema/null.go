package schema

type NullField struct {
	PrimitiveField
}

func NewNullField(definition interface{}) *NullField {
	return &NullField{PrimitiveField{
		definition:       definition,
		name:             "Null",
		goType:           "*types.NullVal",
		serializerMethod: "vm.WriteNull",
	}}
}

func (s *NullField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	return "", nil
}

func (s *NullField) WrapperType() string {
	return ""
}

func (s *NullField) IsReadableBy(f AvroType) bool {
	_, ok := f.(*NullField)
	return ok
}
