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
	return canonicalForm(t, make(map[string]interface{}))
}

func canonicalForm(t schema.AvroType, visited map[string]interface{}) interface{} {
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
			members = append(members, canonicalForm(m, visited))
		}
		return members
	case *schema.ArrayField:
		return &CanonicalFields{
			Type:  "array",
			Items: canonicalForm(v.ItemType(), visited),
		}
	case *schema.MapField:
		return &CanonicalFields{
			Type:   "map",
			Values: canonicalForm(v.ItemType(), visited),
		}
	case *schema.Reference:
		name := v.Def.AvroName().String()
		if _, ok := visited[name]; ok {
			return name
		} else {
			visited[name] = true
			switch def := v.Def.(type) {
			case *schema.RecordDefinition:
				fields := make([]interface{}, 0)
				for _, f := range def.Fields() {
					fn := f.Name()
					fields = append(fields, &CanonicalFields{
						Name: &fn,
						Type: canonicalForm(f.Type(), visited),
					})
				}

				return &CanonicalFields{
					Name:   &name,
					Type:   "record",
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
	}
	panic(fmt.Sprintf("Unkonwn type: %T", t))
}
