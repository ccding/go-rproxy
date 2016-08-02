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

package certs

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

// LoadCACerts loads the CA certificates (root certs).
func LoadCACerts(certsFilename string) (*x509.CertPool, error) {
	certs, err := ioutil.ReadFile(certsFilename)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if ok := cp.AppendCertsFromPEM([]byte(certs)); !ok {
		return nil, errors.New("fail loading certificates")
	}
	return cp, nil
}

// LoadClientCerts loads the client certificates.
func LoadClientCerts(rootCert, clientCert, clientKey, serverName string) (*tls.Config, error) {
	// Load root certificate
	roots, err := LoadCACerts(rootCert)
	if err != nil {
		return nil, err
	}
	// Load client certificate
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}
	// Set TLS config
	config := &tls.Config{
		RootCAs:      roots,
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
	}
	return config, nil
}

// LoadServerCerts loads the server certificates.
func LoadServerCerts(rootCert, serverCert, serverKey string) (*tls.Config, error) {
	// Load root certificate
	roots, err := LoadCACerts(rootCert)
	if err != nil {
		return nil, err
	}
	// Load server certificate
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}
	// Set TLS config
	config := &tls.Config{
		ClientCAs:    roots,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}
	return config, nil
}
