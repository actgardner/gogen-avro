package canonical

import (
	"fmt"

	"github.com/actgardner/gogen-avro/schema"
)

type CanonicalFields struct {
	Name    *string       `json:"name,omitempty"`
	Type    interface{}   `json:"type,omitempty"`
	Fields  []interface{} `json:"fields,omitempty"`
	Symbols []string      `json:"symbols,omitempty"`
	Items   interface{}   `json:"items,omitempty"`
	Values  interface{}   `json:"values,omitempty"`
	Size    *int          `json:"size,omitempty"`
}

func CanonicalForm(t schema.AvroType) interface{} {
	switch v := t.(type) {
	case *schema.BoolField:
		return "boolean"
	case *schema.BytesField:
		return "bytes"
	case *schema.DoubleField:
		return "double"
	case *schema.FloatField:
		return "float"
	case *schema.IntField:
		return "int"
	case *schema.LongField:
		return "long"
	case *schema.NullField:
		return "null"
	case *schema.StringField:
		return "string"
	case *schema.UnionField:
		members := make([]interface{}, 0)
		for _, m := range v.AvroTypes() {
			members = append(members, CanonicalForm(m))
		}
		return members
	case *schema.ArrayField:
		return &CanonicalFields{
			Type:  "array",
			Items: CanonicalForm(v.ItemType()),
		}
	case *schema.MapField:
		return &CanonicalFields{
			Type:   "map",
			Values: CanonicalForm(v.ItemType()),
		}
	case *schema.Reference:
		name := v.Def.AvroName().String()
		switch def := v.Def.(type) {
		case *schema.RecordDefinition:
			fields := make([]interface{}, 0)
			for _, f := range def.Fields() {
				fields = append(fields, &CanonicalFields{
					Name: &name,
					Type: CanonicalForm(f.Type()),
				})
			}

			return &CanonicalFields{
				Name:   &name,
				Fields: fields,
			}
		case *schema.EnumDefinition:
			return &CanonicalFields{
				Name:    &name,
				Type:    "enum",
				Symbols: def.Symbols(),
			}
		case *schema.FixedDefinition:
			size := def.SizeBytes()
			return &CanonicalFields{
				Name: &name,
				Type: "fixed",
				Size: &size,
			}
		}
	}
	panic(fmt.Sprintf("Unkonwn type: %T", t))
}
