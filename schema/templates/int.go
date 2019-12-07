package templates

type IntContext struct {
	PrimitiveContext
}

func NewIntContext() *BytesField {
	return &IntContext{PrimitiveContext{
		goType:           "int32",
		serializerMethod: "vm.WriteInt",
		wrapperType:      "types.Int",
	}}
}

func (_ *IntContext) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}
