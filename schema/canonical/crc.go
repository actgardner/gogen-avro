package canonical

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

// IMplementation of Avro's 64-bit Rabin code for fingerprinting schemas

/*
long fingerprint64(byte[] buf) {
  if (FP_TABLE == null) initFPTable();
  long fp = EMPTY;
  for (int i = 0; i < buf.length; i++)
    fp = (fp >>> 8) ^ FP_TABLE[(int)(fp ^ buf[i]) & 0xff];
  return fp;
}

static long EMPTY = 0xc15d213aa4d7a795L;
static long[] FP_TABLE = null;

void initFPTable() {
  FP_TABLE = new long[256];
  for (int i = 0; i < 256; i++) {
    long fp = i;
    for (int j = 0; j < 8; j++)
      fp = (fp >>> 1) ^ (EMPTY & -(fp & 1L));
    FP_TABLE[i] = fp;
  }
}
*/

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

func AvroCRC64Fingerprint(data []byte) string {
	fp := EMPTY
	for _, d := range data {
		fp = (fp >> 8) ^ FP_TABLE[(fp^uint64(d))&0xff]
	}
	output := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(output, binary.LittleEndian, fp)
	return hex.EncodeToString(output.Bytes())
}
