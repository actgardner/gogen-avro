package avro

import (
	"io"
)

type Codec string

const (
	Null    Codec = "null"
	Deflate Codec = "deflate"
	Snappy  Codec = "snappy"
)

type FlushableResettableWriter interface {
	Flush() error
	Reset(io.Writer)
}
