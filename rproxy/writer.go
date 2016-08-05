package rproxy

import (
	"io"
)

type RPWriteCloser struct {
	Writer io.Writer
	Closer io.Closer
}

func (r *RPWriteCloser) Write(p []byte) (int, error) {
	return r.Writer.Write(p)
}

func (r *RPWriteCloser) Close() error {
	return r.Closer.Close()
}

func NewRPWriteCloser(wc io.WriteCloser) io.WriteCloser {
	return &RPWriteCloser{Writer: wc, Closer: wc}
}
