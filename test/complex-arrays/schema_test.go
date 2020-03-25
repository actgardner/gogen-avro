package avro

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Round trip some records nested in arrays
const fixtureJson = `
[
  {"Children": []},
  {"Children": [{"Name": "test-record"}]},
  {"Children": [{"Name": "test-record"}, {"Name": "test-record-2"}]}
]
`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]Parent, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeParent(&buf)
		assert.Nil(t, err)
		assert.Equal(t, *datum, f)
	}
}
