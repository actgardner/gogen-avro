package schema

type UnionField struct {
	itemType   []AvroType
	definition []interface{}
}

func NewUnionField(itemType []AvroType, definition []interface{}) *UnionField {
	return &UnionField{
		itemType:   itemType,
		definition: definition,
	}
}

func (s *UnionField) ItemTypes() []AvroType {
	return s.itemType
}

func (s *UnionField) Children() []AvroType {
	return s.itemType
}

func (s *UnionField) IsReadableBy(f AvroType) bool {
	// Report if *any* writer type could be deserialized by the reader
	for _, t := range s.ItemTypes() {
		if readerUnion, ok := f.(*UnionField); ok {
			for _, rt := range readerUnion.ItemTypes() {
				if t.IsReadableBy(rt) {
					return true
				}
			}
		} else {
			if t.IsReadableBy(f) {
				return true
			}
		}
	}
	return false
}

func (s *UnionField) Equals(reader *UnionField) bool {
	if len(reader.ItemTypes()) != len(s.ItemTypes()) {
		return false
	}

	for i, t := range s.ItemTypes() {
		readerType := reader.ItemTypes()[i]
		if writerRef, ok := t.(*Reference); ok {
			if readerRef, ok := readerType.(*Reference); ok {
				if readerRef.TypeName != writerRef.TypeName {
					return false
				}
			} else {
				return false
			}
		} else if t != readerType {
			return false
		}
	}
	return true
}
