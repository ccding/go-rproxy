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

package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
)

const (
	appName    = "testapp"
	rootCert   = "certs/root_cert.pem"
	clientCert = "certs/client_0_cert.pem"
	clientKey  = "certs/client_0_key.pem"
	listenAddr = "127.0.0.1:23001" // listen addr of the rproxy
)

func main() {
	// Load root certificate to verify server certificate
	rootPEM, err := ioutil.ReadFile(rootCert)
	if err != nil {
		log.Fatalf("failed to read root certificate: %s", err)
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		log.Fatalf("failed to parse root certificate")
	}
	// Load client certificate
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatalf("failed to load client tls certificate: %s", err)
	}
	// Set TLS config
	config := tls.Config{
		RootCAs:      roots,
		ServerName:   appName + "-server",
		Certificates: []tls.Certificate{cert},
	}
	// Listen to the TLS port
	conn, err := tls.Dial("tcp", listenAddr, &config)
	if err != nil {
		log.Fatalf("error: dial: %s", err)
	}
	defer conn.Close()
	// Read each line from stdin and send it to the server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		log.Printf("sending: %s", message)
		// Send message to the server
		n, err := conn.Write([]byte(message))
		if err != nil {
			log.Fatalf("error: write: %s", err)
		}
		log.Printf("wrote %d bytes", n)
		// Receive response from the server
		buf := make([]byte, 1024)
		n, err = conn.Read(buf)
		if err != nil {
			log.Fatalf("error: read: %s", err)
		}
		log.Printf("read %d bytes", n)
		log.Printf("server reply: %s", string(buf[:n]))
	}
}
