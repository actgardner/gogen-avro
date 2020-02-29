package avro

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Round trip some records nested in arrays
const fixtureJson = `
[
  {"Children": null},
  {"Children": [{"Name": "test-record"}]},
  {"Children": [{"Name": "test-record"}, {"Name": "test-record-2"}]}
]
`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]Parent, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, err := DeserializeParent(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, datum, f)
	}
}
