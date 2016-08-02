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
	"flag"
	"log"
	"strings"

	"github.com/ccding/go-rproxy/rproxy"
)

func main() {
	var listen = flag.String("l", "tls://:23001", "listen address")
	var backend = flag.String("b", "tls://127.0.0.1:23002", "backend address")
	var rootCert = flag.String("rcert", "certs/root_cert.pem", "root cert")
	var serverCert = flag.String("scert", "certs/server_cert.pem", "server cert")
	var serverKey = flag.String("skey", "certs/server_key.pem", "server key")
	var clientCert = flag.String("ccert", "certs/client_0_cert.pem", "client cert")
	var clientKey = flag.String("ckey", "certs/client_0_key.pem", "client key")
	flag.Parse()

	listenProtoAndAddr := strings.Split(*listen, "://")
	backendProtoAndAddr := strings.Split(*backend, "://")

	if len(listenProtoAndAddr) != 2 || len(backendProtoAndAddr) != 2 {
		panic("Wrong arguments.")
	}

	rp := rproxy.NewRProxy(
		listenProtoAndAddr[0],
		listenProtoAndAddr[1],
		backendProtoAndAddr[0],
		backendProtoAndAddr[1],
		*rootCert,
		*serverCert,
		*serverKey,
		*clientCert,
		*clientKey,
	)

	go rp.Start()

	log.Printf("Listening on: %s\n", *listen)
	log.Printf("Forwarding to: %s\n", *backend)

	select {}
}
