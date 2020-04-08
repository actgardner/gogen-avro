package avro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/actgardner/gogen-avro/compiler"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
  {"IntField": [
    {"UnionType": 1, "Int": 1}, 
    null, 
    {"UnionType": 1, "Int": 2147483647}
  ]}
]
`

func TestRoundTrip(t *testing.T) {
	compiler.LoggingEnabled = true
	fixtures := make([]*ArrayTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)
		fmt.Printf("%v\n", buf)

		datum, err := DeserializeArrayTestRecord(&buf)
		assert.Nil(t, err)

		assert.Equal(t, datum, f)
	}
}
