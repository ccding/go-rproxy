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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
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
	serverName := appName + "-server"
	config, err := loadClientCert(serverName, rootCert, clientCert, clientKey)
	if err != nil {
		log.Fatal("%s", err)
	}
	conn, err := tls.Dial("tcp", listenAddr, config)
	if err != nil {
		log.Fatalf("error: dial: %s", err)
	}
	defer conn.Close()
	go read(conn)
	write(conn)
}

func loadClientCert(serverName, rootCert, clientCert, clientKey string) (*tls.Config, error) {
	// Load root certificate to verify server certificate
	rootPEM, err := ioutil.ReadFile(rootCert)
	if err != nil {
		return nil, err
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		return nil, errors.New("failed to parse root certificate")
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

func write(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Fatalf("error: write: %s", err)
		}
	}
}

func read(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatalf("error: read: %s", err)
			}
			os.Exit(0)
		}
		fmt.Println(string(buf[:n]))
	}
}
