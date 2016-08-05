// Copyright 2016, Cong Ding. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: Cong Ding <dinggnu@gmail.com>

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
