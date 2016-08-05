package rproxy

import (
	"io"
)

// RPWriteCloser defines a customized WriteCloser, which is used to modify data.
type RPWriteCloser struct {
	Writer io.Writer
	Closer io.Closer
}

// NewRPWriteCloser creates the new RPWriteCloser from an io.WriteCloser.
func NewRPWriteCloser(wc io.WriteCloser) io.WriteCloser {
	return &RPWriteCloser{Writer: wc, Closer: wc}
}

// Write writes data to the writer.
func (r *RPWriteCloser) Write(p []byte) (int, error) {
	return r.Writer.Write(p)
}

// Close closes the connection.
func (r *RPWriteCloser) Close() error {
	return r.Closer.Close()
}
