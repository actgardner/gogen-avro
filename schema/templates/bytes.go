package templates

type BytesContext struct {
	PrimitiveContext
}

func NewBytesContext() *BytesField {
	return &BytesContext{PrimitiveContext{
		goType:           "[]byte",
		serializerMethod: "vm.WriteBytes",
		wrapperType:      "types.Bytes",
	}}
}

func (_ *BytesContext) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = []byte(%q)", lvalue, rvalue), nil
}
