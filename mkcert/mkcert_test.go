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
// Author: Cong Ding <dinggnu@gmail.com>

package mkcert

import (
	"testing"

	"github.com/ccding/go-rproxy/certs"
)

func TestCACerts(t *testing.T) {
	cp, err := certs.LoadCACerts("root_cert.pem")
	if err != nil {
		t.Errorf("error: failed to load the certificates.")
	}
	if cp == nil {
		t.Errorf("nil cert: failed to load the certificates.")
	}
}
