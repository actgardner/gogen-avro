package avro

import (
	"bytes"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
var fixtures = []Event{
	{
		Id: "id1",
	},
	{
		Id: "differentid",
	},
}

func compareFixtureGoAvro(t *testing.T, actual interface{}, expected interface{}) {
	record := actual.(map[string]interface{})
	fixture := expected.(Event)
	id, ok := record["id"]
	assert.Equal(t, ok, true)
	assert.Equal(t, id, fixture.Id)
}

func TestRootUnionFixture(t *testing.T) {
	codec, err := goavro.NewCodec(fixtures[0].Schema())
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()

		err = writeEvent(&f, &buf)
		assert.Nil(t, err)

		datum, remaining, err := codec.NativeFromBinary(buf.Bytes())
		assert.Nil(t, err)
		assert.Equal(t, 0, len(remaining))

		compareFixtureGoAvro(t, datum, f)
	}
}

func TestRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()

		err := writeEvent(&f, &buf)
		assert.Nil(t, err)

		datum, err := DeserializeEvent(&buf)
		assert.Nil(t, err)
		assert.Equal(t, datum, &f)
	}
}
