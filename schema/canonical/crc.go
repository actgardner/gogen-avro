package canonical

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/actgardner/gogen-avro/parser"
	"github.com/actgardner/gogen-avro/resolver"
	"io"
)

var FP_TABLE []uint64

const EMPTY uint64 = 0xc15d213aa4d7a795

var HeaderV1 = []byte{0xC3, 0x01}

func init() {
	FP_TABLE = make([]uint64, 256)
	for i := range FP_TABLE {
		fp := uint64(i)
		for j := 0; j < 8; j++ {
			fp = (fp >> 1) ^ (EMPTY & -(fp & 1))
		}
		FP_TABLE[i] = fp
	}
}

func AvroCRC64Fingerprint(schema []byte, writer io.Writer) error {
	fp := EMPTY
	for _, d := range schema {
		fp = (fp >> 8) ^ FP_TABLE[(fp^uint64(d))&0xff]
	}
	err := binary.Write(writer, binary.LittleEndian, fp)
	return err
}

func AvroCalcSchemaUID(schema string) []byte {
	b := make([]byte, 0, 8)
	output := bytes.NewBuffer(b)
	ns := parser.NewNamespace(false)

	s, err := ns.TypeForSchema([]byte(schema))
	if err != nil {
		panic("Invalid schema 1")
	}
	for _, def := range ns.Roots {
		err = resolver.ResolveDefinition(def, ns.Definitions)
		if err != nil {
			panic("Invalid schema 2")
		}
	}
	canonical, err := json.Marshal(CanonicalForm(s))
	err = AvroCRC64Fingerprint(canonical, output)
	if err != nil {
		panic("Invalid schema 1")
	}
	return output.Bytes()
}

func AvroVersionHeader(writer io.Writer, header []byte) error {
	fp := HeaderV1
	err := binary.Write(writer, binary.LittleEndian, fp)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, header)
	return err
}

func AvroConsumeHeader(r io.Reader) error {
	buf := make([]byte, 10)
	n, err := r.Read(buf)
	if n < 10 || err != nil {
		return errors.New("message length <2")
	}
	return nil
}
