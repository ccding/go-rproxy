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
	appName     = "testapp"
	rootCert    = "certs/root_cert.pem"
	serverCert  = "certs/server_cert.pem"
	serverKey   = "certs/server_key.pem"
	backendAddr = "127.0.0.1:23002" // backend addr of the rproxy
)

func main() {
	config, err := certs.LoadServerCerts(rootCert, serverCert, serverKey)
	if err != nil {
		log.Fatalf("load server certs error: %v", err)
	}
	listener, err := tls.Listen("tcp", backendAddr, config)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	var stdin = make(chan string, 1024)
	go readStdin(stdin)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("accept error: %v", err)
		}
		go handleConn(conn, stdin)
	}
}

func readStdin(stdin chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		stdin <- scanner.Text() + "\n"
	}
}

func handleConn(conn net.Conn, stdin chan string) {
	defer conn.Close()
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Fatalf("error: not tls conn")
		return
	}
	if err := tlsConn.Handshake(); err != nil {
		log.Fatalf("handshake error: %v", err)
		return
	}
	var quit = make(chan bool, 2)
	go read(conn, quit)
	write(conn, quit, stdin)
}

func read(conn net.Conn, quit chan bool) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatalf("read error: %v", err)
			}
			quit <- true
			return
		}
		fmt.Print(string(buf[:n]))
	}
}

func write(conn net.Conn, quit chan bool, stdin chan string) {
	for {
		select {
		case <-quit:
			return
		case message := <-stdin:
			_, err := conn.Write([]byte(message))
			if err != nil {
				quit <- true
				return
			}
		}
	}
}
