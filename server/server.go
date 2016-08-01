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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

const (
	appName     = "testapp"
	rootCert    = "certs/root_cert.pem"
	serverCert  = "certs/server_cert.pem"
	serverKey   = "certs/server_key.pem"
	backendAddr = "127.0.0.1:23002" // backend addr of the rproxy
)

func main() {
	// Load root certificate to verify client certificate
	rootPEM, err := ioutil.ReadFile(rootCert)
	if err != nil {
		log.Fatalf("failed to read root certificate: %s", err)
	}
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM([]byte(rootPEM)); !ok {
		log.Fatalf("failed to parse root certificate")
	}
	// Load server certificate
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatalf("failed to load server tls certificate: %s", err)
	}
	// Set TLS config
	config := tls.Config{
		ClientCAs:    roots,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}
	// Listen to the TLS port
	listener, err := tls.Listen("tcp", backendAddr, &config)
	if err != nil {
		log.Fatalf("error: listen: %s", err)
	}
	log.Printf("server started")
	// Handle requests
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: accept: %s", err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Fatalf("error: not tls conn")
		return
	}

	err := tlsConn.Handshake()
	if err != nil {
		log.Fatalf("error: handshake: %s", err)
		return
	}

	clientID, err := getClientID(tlsConn)
	if err != nil {
		log.Printf("error: cannot get client-id: %s", err)
		return
	}
	log.Printf("handle connection from client-id: %s", clientID)
	// Handle client message
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("client:%s: error: read: %s", clientID, err)
			}
			break
		}
		log.Printf("client:%s: read %d bytes", clientID, n)
		log.Printf("client:%s: echo: %s", clientID, string(buf[:n]))
		// Echo message back to the client
		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("client:%s: error: write: %s", clientID, err)
			break
		}
		log.Printf("client:%s: wrote %d bytes", clientID, n)
	}
	log.Printf("client:%s: connection closed", clientID)
}

func getClientID(tlsConn *tls.Conn) (string, error) {
	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return "", fmt.Errorf("client certificate not found")
	}

	cert := state.PeerCertificates[0]
	parts := strings.Split(cert.Subject.CommonName, "-")
	if len(parts) != 3 || parts[0] != appName || parts[1] != "client" || len(parts[2]) == 0 {
		return "", fmt.Errorf("bad client common name")
	}

	return parts[2], nil
}
