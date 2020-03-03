package schema

type NullField struct {
	primitiveField
}

func NewNullField(definition interface{}) *NullField {
	return &NullField{newPrimitiveField("Null", "types.NullVal", definition, "vm.WriteNull")}
}

func (s *NullField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	if _, ok := f.(*NullField); ok {
		return true
	}
	if s.primitiveField.IsReadableBy(f, visited) {
		return true
	}
	if union, ok := f.(*UnionField); ok {
		return isReadableByUnion(s, union, visited)
	}
	return false
}
