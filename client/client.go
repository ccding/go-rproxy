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
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/ccding/go-rproxy/certs"
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
	config, err := certs.LoadClientCerts(rootCert, clientCert, clientKey, serverName)
	if err != nil {
		log.Fatalf("%s", err)
	}
	conn, err := tls.Dial("tcp", listenAddr, config)
	if err != nil {
		log.Fatalf("error: dial: %s", err)
	}
	defer conn.Close()
	go read(conn)
	write(conn)
}

func write(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text() + "\n"
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
		fmt.Print(string(buf[:n]))
	}
}
