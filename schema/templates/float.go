package templates

type FloatContext struct {
	PrimitiveContext
}

func NewFloatContext() *BytesField {
	return &FloatContext{PrimitiveContext{
		goType:           "float32",
		serializerMethod: "vm.WriteFloat",
		wrapperType:      "types.Float",
	}}
}

func (_ *FloatContext) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected float as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}
