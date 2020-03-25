package avro

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
{"RecursiveField":{"UnionType":0}},
{"RecursiveField":{"RecursiveUnionTestRecord":{"RecursiveField": {"UnionType": 0}}, "UnionType":1}}
]
`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]RecursiveUnionTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeRecursiveUnionTestRecord(&buf)
		assert.Nil(t, err)

		assert.Equal(t, *datum, f)
	}
}
