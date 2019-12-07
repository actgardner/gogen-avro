package templates

import (
	"github.com/actgardner/gogen-avro/schema"
)

type Context interface {
	GoType() string
	ConstructorMethod() string
	SerializerMethod() string
	Template() string
	DefaultValue(lvalue string, rvalue interface{}) (string, error)
	WrapperType() string
}

func ContextForField(field schema.Field) Context {
	switch fType := field.(type) {
	case *schema.Reference:
		return ContextForReference(fType)
	case *schema.ArrayField:
		return NewArrayContext(fType)
	case *schema.MapField:
		return NewMapContext(fType)
	case *schema.BoolField:
		return NewBoolContext()
	case *schema.BytesField:
		return NewBytesContext()
	case *schema.DoubleField:
		return NewDoubleContext()
	case *schema.FloatField:
		return NewFloatContext()
	case *schema.IntField:
		return NewIntContext()
	case *schema.LongField:
		return NewIntContext()
	}
	panic("Unsupported field type: %v", field)
}

func ContextForReference(ref *schema.Reference) Context {
	switch def := ref.Def.(type) {
	case *schema.RecordDefinition:
		return NewRecordContext(def)
	}
	panic("Unsupported reference type: %v", ref)
}
