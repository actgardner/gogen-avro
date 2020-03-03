package schema

// Common methods for all primitive types
type primitiveField struct {
	abstractField
	serializerMethod string
}

func newPrimitiveField(name, goType string, definition interface{}, serializerMethod string) primitiveField {
	return primitiveField{
		abstractField: abstractField{
			name:       name,
			simpleName: name,
			goType:     goType,
			definition: definition,
		},
		serializerMethod: serializerMethod,
	}
}

var (
	// Ensure interface implementation
	_ AvroType = &primitiveField{}
)

func (p *primitiveField) SerializerMethod() string {
	return p.serializerMethod
}

func (p *primitiveField) setSerializerMethod(serializerMethod string) {
	p.serializerMethod = serializerMethod
}

func (p *primitiveField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	return "", nil
}

func (p *primitiveField) WrapperType() string {
	return ""
}

func (p *primitiveField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	if p.goType == f.GoType() {
		return true
	}
	return false
}

func isReadableByUnion(t AvroType, union *UnionField, visited map[QualifiedName]interface{}) bool {
	for _, f := range union.Children() {
		if t.IsReadableBy(f, visited) {
			return true
		}
	}
	return false
}
