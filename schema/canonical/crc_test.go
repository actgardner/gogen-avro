package canonical

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/actgardner/gogen-avro/parser"
	"github.com/actgardner/gogen-avro/resolver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAvroCRC64Fingerprint(t *testing.T) {
	cases := []struct {
		schema      string
		fingerprint string
	}{
		{"{\"type\":\"record\",\"name\":\"SmsMessage\",\"namespace\":\"ca.bell.vas.crossfunctional.domain\",\"fields\":[{\"name\":\"versionNumber\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"},{\"name\":\"transactionType\",\"type\":{\"type\":\"enum\",\"name\":\"TransactionType\",\"symbols\":[\"MSG\",\"ACK\",\"RSP\"]},\"default\":\"MSG\"},{\"name\":\"orignalShortMessageInMsgPayload\",\"type\":\"boolean\",\"default\":false},{\"name\":\"startTimestamp\",\"type\":\"long\",\"default\":0},{\"name\":\"internalTransactionId\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"},{\"name\":\"externalTransactionId\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"},{\"name\":\"queueResponseName\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"},{\"name\":\"responseSequenceNumber\",\"type\":\"int\",\"default\":0},{\"name\":\"messageSource\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"},{\"name\":\"messageDestination\",\"type\":{\"type\":\"enum\",\"name\":\"TransactionDestination\",\"symbols\":[\"BELL\",\"LONGCODE\",\"ICMS\",\"ROGERS\",\"SYNIVERSE\",\"EXT\",\"UNKNOWN\"]},\"default\":\"UNKNOWN\"},{\"name\":\"messageDestinationCode\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"},{\"name\":\"segmentSequenceNumber\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Short\"},\"default\":0},{\"name\":\"totalNumberOfSegment\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Short\"},\"default\":0},{\"name\":\"segmentMessageRefNumber\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Short\"},\"default\":0},{\"name\":\"contentFilteringResult\",\"type\":{\"type\":\"enum\",\"name\":\"ContentFilteringResult\",\"symbols\":[\"BLOCKED\",\"ALLOWED\",\"NULL\"]},\"default\":\"ALLOWED\"},{\"name\":\"pdu\",\"type\":{\"type\":\"record\",\"name\":\"Pdu\",\"fields\":[{\"name\":\"name\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"},{\"name\":\"isRequest\",\"type\":\"boolean\",\"default\":false},{\"name\":\"commandLength\",\"type\":\"int\",\"default\":0},{\"name\":\"commandId\",\"type\":\"int\",\"default\":4},{\"name\":\"commandStatus\",\"type\":\"int\",\"default\":0},{\"name\":\"sequenceNumber\",\"type\":\"int\",\"default\":0},{\"name\":\"serviceType\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"},{\"name\":\"sourceAddress\",\"type\":{\"type\":\"record\",\"name\":\"Address\",\"fields\":[{\"name\":\"ton\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Byte\"},\"default\":0},{\"name\":\"npi\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Byte\"},\"default\":0},{\"name\":\"address\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"}]},\"default\":{}},{\"name\":\"destAddress\",\"type\":\"Address\",\"default\":{}},{\"name\":\"esmClass\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Byte\"},\"default\":0},{\"name\":\"protocolId\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Byte\"},\"default\":0},{\"name\":\"priority\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Byte\"},\"default\":0},{\"name\":\"scheduleDeliveryTime\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"},{\"name\":\"validityPeriod\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"},{\"name\":\"registeredDelivery\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Byte\"},\"default\":0},{\"name\":\"replaceIfPresent\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Byte\"},\"default\":0},{\"name\":\"dataCoding\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Byte\"},\"default\":0},{\"name\":\"defaultMsgId\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Byte\"},\"default\":0},{\"name\":\"shortMessage\",\"type\":{\"type\":\"bytes\",\"java-class\":\"[B\"},\"default\":\"ÿ\"},{\"name\":\"optionalParameters\",\"type\":{\"type\":\"array\",\"items\":{\"type\":\"record\",\"name\":\"Tlv\",\"fields\":[{\"name\":\"tag\",\"type\":{\"type\":\"int\",\"java-class\":\"java.lang.Short\"},\"default\":0},{\"name\":\"value\",\"type\":{\"type\":\"bytes\",\"java-class\":\"[B\"},\"default\":\"ÿ\"},{\"name\":\"tagName\",\"type\":{\"type\":\"string\",\"avro.java.string\":\"String\"},\"default\":\"\"}]},\"java-class\":\"java.util.ArrayList\"},\"default\":[]}]},\"default\":{}}],\"default\":{}}", "0c947f601de7ce84"},
	}

	for _, c := range cases {
		b := make([]byte, 0, 8)
		output := bytes.NewBuffer(b)
		ns := parser.NewNamespace(false)
		s, err := ns.TypeForSchema([]byte(c.schema))
		assert.Nil(t, err)
		for _, def := range ns.Roots {
			assert.Nil(t, resolver.ResolveDefinition(def, ns.Definitions))
		}
		canonical, err := json.Marshal(CanonicalForm(s))
		err = AvroCRC64Fingerprint(canonical, output)
		assert.Nil(t, err)
		assert.Equal(t, c.fingerprint, hex.EncodeToString(output.Bytes()))
	}
}

func TestAvroVersionHeader(t *testing.T) {
	cases := []struct {
		header   []byte
		expected string
	}{
		{[]byte{0x0c, 0x94, 0x7f, 0x60, 0x1d, 0xe7, 0xce, 0x84}, "c3010c947f601de7ce84"},
	}

	for _, c := range cases {
		b := make([]byte, 0, 2)
		output := bytes.NewBuffer(b)
		err := AvroVersionHeader(output, c.header)
		assert.Nil(t, err)
		assert.Equal(t, c.expected, hex.EncodeToString(output.Bytes()))
	}
}
