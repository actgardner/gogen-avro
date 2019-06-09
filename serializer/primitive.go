package serializer

// Reader is the interface that wraps the basic Read method.
//
// ReadNext reads the next block of bytes from the underlaying data stream.
// It returns the bytes read inside a buffer and any error encountered.
//
// Implementations define the type of the returned byte buffer.
type Reader interface {
	ReadNext() ([]byte, error)
}

// Writer is the interface that wraps the basic Write method.
//
// Write writes len(p) bytes from p to the underlying data stream.
// A expected length header of len(p) is written before the actual message.
// It returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p). Write must not modify the slice data, even temporarily.
//
// Implementations define the type of p.
type Writer interface {
	Write(p []byte) error
}
