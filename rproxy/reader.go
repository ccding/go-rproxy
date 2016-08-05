package rproxy

import (
	"io"
)

// RPReader defines a customized Reader, which is used to modify or log data.
type RPReader struct {
	Reader io.Reader
}

// NewRPReader creates the new RPReader from an io.Reader
func NewRPReader(r io.Reader) io.Reader {
	return &RPReader{Reader: r}
}

// Read reads data from the reader.
func (r *RPReader) Read(p []byte) (int, error) {
	return r.Reader.Read(p)
}
