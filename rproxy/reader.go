package rproxy

import (
	"io"
	"log"
)

// RPReader defines a customized Reader, which is used to log data.
type RPReader struct {
	Reader  io.Reader
	verbose bool
}

// NewRPReader creates the new RPReader from an io.Reader.
func NewRPReader(r io.Reader, v bool) io.Reader {
	return &RPReader{Reader: r, verbose: v}
}

// Read reads data from the reader.
func (r *RPReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if r.verbose {
		log.Print(string(p[:n]))
	}
	return
}
