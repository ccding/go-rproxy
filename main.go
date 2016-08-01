package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ccding/go-rproxy/rproxy"
)

func main() {
	var listen = flag.String("l", "tcp://:23001", "listen address")
	var backend = flag.String("b", "tcp://127.0.0.1:23002", "backend address")
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

	fmt.Printf("Listening on: %s\n", *listen)
	fmt.Printf("Forwarding to: %s\n", *backend)

	select {}
}
