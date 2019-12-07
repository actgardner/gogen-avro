package templates

type StringContext struct {
	PrimitiveContext
}

func NewStringContext() *BytesField {
	return &StringContext{PrimitiveContext{
		goType:           "string",
		serializerMethod: "vm.WriteString",
		wrapperType:      "types.String",
	}}
}

func (_ *StringContext) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %q", lvalue, rvalue), nil
}
