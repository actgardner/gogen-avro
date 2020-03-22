package avro

import (
	"bytes"
	"encoding/json"
	"github.com/actgardner/gogen-avro/singleobject"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
  {
    "OtherField": {
      "StringField": "789",
      "BoolField": true,
      "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"
    }
  },
  {
    "OtherField": {
      "StringField": "abcdghejw",
      "BoolField": true,
      "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"
    }
  },
  {
    "OtherField": {
      "StringField": "jdnwjkendwedddedee",
      "BoolField": true,
      "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"
    }
  }
]
`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]NestedTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		writer := singleobject.NewWriter(&buf, NestedRecordUID)
		err = f.Serialize(writer)
		assert.Nil(t, err)

		datum, err := DeserializeNestedTestRecord(singleobject.NewReader(&buf))
		assert.Nil(t, err)
		assert.Equal(t, *datum, f)
	}
}
