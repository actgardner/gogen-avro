package templates

type NullContext struct {
	PrimitiveContext
}

func NewIntContext() *BytesField {
	return &IntContext{PrimitiveContext{
		goType:           "*types.NullVal",
		serializerMethod: "vm.WriteNull",
		wrapperType:      "",
	}}
}

func (s *NullField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	return "", nil
}
