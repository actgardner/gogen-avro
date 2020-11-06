package avro

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	var fixedDefault = NewFixedDefaultTestRecord()
	fixedDefault.SetDefault(0)

	expected := FixedDefaultTestRecord{
		FixedField: [12]byte{0, 1, 18, 0, 19, 67, 0, 1, 18, 0, 19, 83},
	}
	assert.Equal(t, &expected, fixedDefault, "Comparing default value")
}
