# go-rproxy

[![Build Status](https://travis-ci.org/ccding/go-rproxy.svg?branch=master)](https://travis-ci.org/ccding/go-rproxy)
[![License](https://img.shields.io/badge/License-Apache%202.0-red.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/ccding/go-rproxy?status.svg)](http://godoc.org/github.com/ccding/go-rproxy/rproxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/ccding/go-rproxy)](https://goreportcard.com/report/github.com/ccding/go-rproxy)

WARNING: This project is still under development. Its API may change
significantly.

go-rproxy is a transport layer reverse proxy speaking TCP and TLS.

A general usecase is that it accepts TLS/HTTPS requests and forwards them to
the backend server as TCP/HTTP requests. This helps software engineers to make
their TCP-based applications have TLS interfaces without modifying their
applications.

Another usecase is that it acts as a transport layer proxy to help debugging
the API between clients and the server.

In general, go-rproxy accepts TCP/TLS requests from clients and sends TCP/TLS
requests to the server, in any combination, say TCP->TCP, TCP->TLS, TLS-TCP,
TLS->TLS.

More details please see `main.go`.

To test this library, you can use these tools to send or receive TCP/TLS
requests:
```
server/server.go:	TLS server
client/client.go:	TLS client
nc (netcat):		TCP server/client
telnet:			TCP client
```

More features are on the way...
