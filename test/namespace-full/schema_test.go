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
{"Header": {"UnionType": 0}, "Body": {"UnionType": 0}},
{"Header": {"UnionType": 1, "HeaderworksData": {"UUID": {"UnionType": 0}, "Hostname": {"UnionType": 0}, "Trace": {"UnionType": 0}}}, "Body": {"UnionType": 0}},
{"Header": {"UnionType": 0}, "Body": {"UnionType": 1, "BodyworksData": {"UUID": {"UnionType": 0}, "Hostname": {"UnionType": 0}, "Trace": {"UnionType": 0}}}},
{"Header": {"UnionType": 1, "HeaderworksData": {"UUID": {"UnionType": 0}, "Hostname": {"UnionType": 0}, "Trace": {"UnionType": 0}}}, "Body": {"UnionType": 1, "BodyworksData": {"UUID": {"UnionType": 0}, "Hostname": {"UnionType": 0}, "Trace": {"UnionType": 0}}}}
]`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]ComAvroTestSample, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		writer := singleobject.NewWriter(&buf, ComAvroTestSampleAvroCRC64Fingerprint)
		err = f.Serialize(writer)
		assert.Nil(t, err)

		datum, err := DeserializeComAvroTestSample(singleobject.NewReader(&buf))
		assert.Nil(t, err)
		assert.Equal(t, *datum, f)
	}
}
