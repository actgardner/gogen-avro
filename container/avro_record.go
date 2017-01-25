package container

import (
	"io"
)

type AvroRecord interface {
	Serialize(io.Writer) error
	Schema() string
}
