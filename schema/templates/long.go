package templates

type LongContext struct {
	PrimitiveContext
}

func NewLongContext() *BytesField {
	return &LongContext{PrimitiveContext{
		goType:           "int64",
		serializerMethod: "vm.WriteLong",
		wrapperType:      "types.Long",
	}}
}

func (_ *LongField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for Field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}
