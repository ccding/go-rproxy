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
