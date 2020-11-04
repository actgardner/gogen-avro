package generic

import (
	"github.com/actgardner/gogen-avro/v7/schema"
	"github.com/actgardner/gogen-avro/v7/vm/types"
)

type Datum interface {
	types.Field
	Datum() interface{}
}

func DatumForType(t schema.AvroType) Datum {
	switch st := t.(type) {
	case *schema.BoolField, *schema.BytesField, *schema.FloatField, *schema.DoubleField, *schema.IntField, *schema.LongField, *schema.StringField, *schema.NullField:
		return &primitiveDatum{}
	case *schema.MapField:
		return &mapDatum{itemType: st.ItemType()}
	case *schema.ArrayField:
		return &arrayDatum{itemType: st.ItemType()}
	case *schema.Reference:
		return datumForReference(st)
	case *schema.UnionField:
		return &unionDatum{itemTypes: st.ItemTypes()}
	}
	panic("")
}

func datumForReference(ref *schema.Reference) Datum {
	switch d := ref.Def.(type) {
	case *schema.RecordDefinition:
		return newRecordDatum(d)
	case *schema.EnumDefinition:
		return newEnumDatum(d)
	case *schema.FixedDefinition:
		return &primitiveDatum{}
	}
	panic("")
}
