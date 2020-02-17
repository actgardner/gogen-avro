package canonical

import (
	"encoding/json"
	"testing"

	"github.com/actgardner/gogen-avro/parser"
	"github.com/actgardner/gogen-avro/resolver"
	"github.com/stretchr/testify/assert"
)

const inputSchema = `{
	"type": "record",
	"name": "TestRecord",
	"namespace": "a",
	"doc": "A test record",
	"fields": [
		{
			"name": "field1",
			"type": {"type": "int"},
			"doc": "field docstring",
			"default": 1234
		},
		{
			"name": "field2",
			"type": {
				"type": "map",
				"values": "string"
			},
			"doc": "field docstring",
			"default": {"a": "2"}
		},
		{
			"name": "field3",
			"type": {
				"type": "array",
				"items": "int"
			},
			"doc": "field docstring",
			"default": [1, 2, 3]
		},
		{
			"name": "field4",
			"type": {
				"type": "enum",
				"name": "b.testenum",
				"symbols": ["x", "y", "z"]
			}
		},
		{
			"name": "field5",
			"type": {
				"type": "fixed",
				"name": "testfixed",
				"namespace": "namespacec",
				"size": 5
			}
		}
	]
}
`

const expected = `{"name":"a.TestRecord","fields":[{"name":"a.TestRecord","type":"int"},{"name":"a.TestRecord","type":{"type":"map","values":"string"}},{"name":"a.TestRecord","type":{"type":"array","items":"int"}},{"name":"a.TestRecord","type":{"name":"b.testenum","type":"enum","symbols":["x","y","z"]}},{"name":"a.TestRecord","type":{"name":"namespacec.testfixed","type":"fixed","size":5}}]}`

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
