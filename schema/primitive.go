package schema

// All primitive types embed this
type PrimitiveField struct {
	definition interface{}
}

func (s *PrimitiveField) Children() []AvroType {
	return []AvroType{}
}

func (s *PrimitiveField) Definition(_ map[QualifiedName]interface{}) (interface{}, error) {
	return s.definition, nil
}

// Types for each of the primitives in the Avro spec

type BoolField struct {
	PrimitiveField
}

func NewBoolField(definition interface{}) *BoolField {
	return &BoolField{PrimitiveField{
		definition: definition,
	}}
}

func (_ *BoolField) IsReadableBy(f AvroType) bool {
	_, ok := f.(*BoolField)
	return ok
}

type BytesField struct {
	PrimitiveField
}

func NewBytesField(definition interface{}) *BytesField {
	return &BytesField{PrimitiveField{
		definition: definition,
	}}
}

func (_ *BytesField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*BytesField); ok {
		return true
	}
	if _, ok := f.(*StringField); ok {
		return true
	}
	return false
}

type DoubleField struct {
	PrimitiveField
}

func NewDoubleField(definition interface{}) *DoubleField {
	return &DoubleField{PrimitiveField{
		definition: definition,
	}}
}

func (_ *DoubleField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}

type FloatField struct {
	PrimitiveField
}

func NewFloatField(definition interface{}) *FloatField {
	return &FloatField{PrimitiveField{
		definition: definition,
	}}
}

func (_ *FloatField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*FloatField); ok {
		return true
	}
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}

type IntField struct {
	PrimitiveField
}

func NewIntField(definition interface{}) *IntField {
	return &IntField{PrimitiveField{
		definition: definition,
	}}
}

func (s *IntField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*IntField); ok {
		return true
	}
	if _, ok := f.(*LongField); ok {
		return true
	}
	if _, ok := f.(*FloatField); ok {
		return true
	}
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}

type LongField struct {
	PrimitiveField
}

func NewLongField(definition interface{}) *LongField {
	return &LongField{PrimitiveField{
		definition: definition,
	}}
}

func (_ *LongField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*LongField); ok {
		return true
	}
	if _, ok := f.(*FloatField); ok {
		return true
	}
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}

type NullField struct {
	PrimitiveField
}

func NewNullField(definition interface{}) *NullField {
	return &NullField{PrimitiveField{
		definition: definition,
	}}
}

func (_ *NullField) IsReadableBy(f AvroType) bool {
	_, ok := f.(*NullField)
	return ok
}

type StringField struct {
	PrimitiveField
}

func NewStringField(definition interface{}) *StringField {
	return &StringField{PrimitiveField{
		definition: definition,
	}}
}

func (_ *StringField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*BytesField); ok {
		return true
	}
	if _, ok := f.(*StringField); ok {
		return true
	}
	return false
}
