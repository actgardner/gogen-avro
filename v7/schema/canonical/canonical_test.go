package canonical

import (
	"encoding/json"
	"testing"

	"github.com/actgardner/gogen-avro/v7/parser"
	"github.com/actgardner/gogen-avro/v7/resolver"
	"github.com/stretchr/testify/assert"
)

const inputSchema = `
{
      "type": "record",
      "name": "Interop",
      "namespace": "org.apache.avro",
      "fields": [
        {"name": "intField", "type": "int"},
        {"name": "longField", "type": "long"},
        {"name": "stringField", "type": "string"},
        {"name": "boolField", "type": "boolean"},
        {"name": "floatField", "type": "float"},
        {"name": "doubleField", "type": "double"},
        {"name": "bytesField", "type": "bytes"},
        {"name": "nullField", "type": "null"},
        {"name": "arrayField", "type": {"type": "array", "items": "double"}},
        {
          "name": "mapField",
          "type": {
            "type": "map",
            "values": {"name": "Foo",
                       "type": "record",
                       "fields": [{"name": "label", "type": "string"}]}
          }
        },
        {
          "name": "unionField",
          "type": ["boolean", "double", {"type": "array", "items": "bytes"}]
        },
        {
          "name": "enumField",
          "type": {"type": "enum", "name": "Kind", "symbols": ["A", "B", "C"]}
        },
        {
          "name": "fixedField",
          "type": {"type": "fixed", "name": "MD5", "size": 16}
        },
        {
          "name": "recordField",
          "type": {"type": "record",
                   "name": "Node",
                   "fields": [{"name": "label", "type": "string"},
                              {"name": "children",
                               "type": {"type": "array",
                                        "items": "Node"}}]}
        }
      ]
    }
`

const expected = `{"name":"org.apache.avro.Interop","type":"record","fields":[{"name":"intField","type":"int"},{"name":"longField","type":"long"},{"name":"stringField","type":"string"},{"name":"boolField","type":"boolean"},{"name":"floatField","type":"float"},{"name":"doubleField","type":"double"},{"name":"bytesField","type":"bytes"},{"name":"nullField","type":"null"},{"name":"arrayField","type":{"type":"array","items":"double"}},{"name":"mapField","type":{"type":"map","values":{"name":"org.apache.avro.Foo","type":"record","fields":[{"name":"label","type":"string"}]}}},{"name":"unionField","type":["boolean","double",{"type":"array","items":"bytes"}]},{"name":"enumField","type":{"name":"org.apache.avro.Kind","type":"enum","symbols":["A","B","C"]}},{"name":"fixedField","type":{"name":"org.apache.avro.MD5","type":"fixed","size":16}},{"name":"recordField","type":{"name":"org.apache.avro.Node","type":"record","fields":[{"name":"label","type":"string"},{"name":"children","type":{"type":"array","items":"org.apache.avro.Node"}}]}}]}`

func TestCanonicalForm(t *testing.T) {
	ns := parser.NewNamespace(false)
	s, err := ns.TypeForSchema([]byte(inputSchema))
	assert.Nil(t, err)
	for _, def := range ns.Roots {
		assert.Nil(t, resolver.ResolveDefinition(def, ns.Definitions))
	}
	canonical, err := json.Marshal(CanonicalForm(s))
	assert.Nil(t, err)
	assert.Equal(t, expected, string(canonical))
}
