package templates

type DoubleContext struct {
	PrimitiveContext
}

func NewDoubleContext() *BytesField {
	return &DoubleContext{PrimitiveContext{
		goType:           "float64",
		serializerMethod: "vm.WriteDouble",
		wrapperType:      "types.Double",
	}}
}

func (_ *DoubleContext) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for field %v, got %q", lvalue, rvalue)
	}
	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}
