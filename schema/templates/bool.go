package templates

type BoolContext struct {
	PrimitiveContext
}

func NewBoolContext() *BoolField {
	return &BoolContext{PrimitiveContext{
		goType:           "bool",
		serializerMethod: "vm.WriteBool",
		wrapperType:      "types.Boolean",
	}}
}

func (_ *BoolContext) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(bool); !ok {
		return "", fmt.Errorf("Expected bool as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}
