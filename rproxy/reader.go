package rproxy

import (
	"io"
)

type RPReader struct {
	Reader io.Reader
}

func NewRPReader(r io.Reader) io.Reader {
	return &RPReader{Reader: r}
}

func (r *RPReader) Read(p []byte) (int, error) {
	return r.Reader.Read(p)
}
