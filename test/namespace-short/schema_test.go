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
{"Header": null, "Body": null},
{"Header": {"Data": {"UUID": null, "Hostname": null, "Trace": null}}, "Body": null},
{"Header": null, "Body": {"Data": {"UUID": null, "Hostname": null, "Trace": null}}},
{"Header": {"Data": {"UUID": null, "Hostname": null, "Trace": null}}, "Body": {"Data": {"UUID": null, "Hostname": null, "Trace": null}}}
]`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]TestSample, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeTestSample(&buf)
		assert.Nil(t, err)
		assert.Equal(t, *datum, f)
	}
}
