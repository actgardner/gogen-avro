package avro

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
  {"IntField": { 
      "a": {"int": 1}, 
      "b": null, 
      "c": {"int": 2147483647}
  } 
  }
]
`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]*MapTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeMapTestRecord(&buf)
		assert.Nil(t, err)

		assert.Equal(t, datum, f)
	}
}
