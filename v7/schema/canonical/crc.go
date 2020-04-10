package canonical

import (
	"bytes"
	"encoding/binary"
)

var FP_TABLE []uint64

const EMPTY uint64 = 0xc15d213aa4d7a795

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

func AvroCRC64Fingerprint(data []byte) []byte {
	fp := EMPTY
	for _, d := range data {
		fp = (fp >> 8) ^ FP_TABLE[(fp^uint64(d))&0xff]
	}
	output := bytes.NewBuffer(make([]byte, 0, 8))
	err := binary.Write(output, binary.LittleEndian, fp)
	if err != nil {
		return nil
	}
	return output.Bytes()
}
