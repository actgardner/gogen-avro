package avro

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `[
	{ "OptField": {"m": {"oneEmptyKey": null}} },
	{ "OptField": {"m": {"name": {"String": "johndoe", "UnionType":1}} } },
	{ "OptField": {"m": {"age": {"Int": 42, "UnionType":2}} } },
	{ "OptField": {"m": {"name": {"String": "johndoe", "UnionType":1}, "age": {"Int": 42, "UnionType":2}} } },
	{ "OptField": {"m": {"name": {"String": "johndoe", "UnionType":1}, "emptyKey": null, "age": {"Int": 42, "UnionType":2}} } }
]`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]MapOptionalTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeMapOptionalTestRecord(&buf)
		assert.Nil(t, err)
		assert.Equal(t, datum.OptField.M, f.OptField.M)
	}
}
