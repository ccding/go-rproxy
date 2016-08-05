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

package rproxy

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/ccding/go-rproxy/certs"
)

// RProxy is used to store configurations of the reverse proxy and start
// running the proxy.
type RProxy struct {
	listenProto  string
	listenAddr   string
	backendProto string
	backendAddr  string
	rootCert     string
	serverCert   string
	serverKey    string
	clientCert   string
	clientKey    string
	clientConfig *tls.Config
	serverConfig *tls.Config
	serverName   string
	verbose      bool
}

// NewRProxyWithoutCerts creates an RProxy instance without setting
// certificate files, which is used in pure TCP setting or set the
// TLS configuration later.
func NewRProxyWithoutCerts(listenProto, listenAddr, backendProto, backendAddr string) *RProxy {
	return &RProxy{
		listenProto:  strings.ToLower(listenProto),
		listenAddr:   strings.ToLower(listenAddr),
		backendProto: strings.ToLower(backendProto),
		backendAddr:  strings.ToLower(backendAddr),
	}
}

// NewRProxy creates an RProxy instance.
func NewRProxy(listenProto, listenAddr, backendProto, backendAddr, rootCert, serverCert, serverKey, clientCert, clientKey, serverName string) *RProxy {
	return &RProxy{
		listenProto:  strings.ToLower(listenProto),
		listenAddr:   strings.ToLower(listenAddr),
		backendProto: strings.ToLower(backendProto),
		backendAddr:  strings.ToLower(backendAddr),
		rootCert:     rootCert,
		serverCert:   serverCert,
		serverKey:    serverKey,
		clientCert:   clientCert,
		clientKey:    clientKey,
		serverName:   serverName,
	}
}

// SetVerbose sets the verbose mode, which prints all the data sent and received.
func (rp *RProxy) SetVerbose(v bool) {
	rp.verbose = v
}

// SetClientConfig sets the config for client (backend TLS).
func (rp *RProxy) SetClientConfig(config *tls.Config) {
	rp.clientConfig = config
}

// SetServerConfig sets the config for server (listen TLS).
func (rp *RProxy) SetServerConfig(config *tls.Config) {
	rp.serverConfig = config
}

// Start starts the reverse proxy service.
func (rp *RProxy) Start() error {
	// Check backend protocol and load certificates if TLS
	switch rp.backendProto {
	case "tcp":
	case "tls":
		// Load client certificates for TLS
		if rp.clientConfig == nil {
			config, err := certs.LoadClientCerts(rp.rootCert, rp.clientCert, rp.clientKey, rp.serverName)
			if err != nil {
				return err
			}
			rp.clientConfig = config
		}
	default:
		return errors.New("backend protocol not supported")
	}
	// Check listen protocol, load certiticates if TLS, and start listening
	switch rp.listenProto {
	case "tcp":
		return rp.startTCP()
	case "tls":
		// Load server certificates for TLS
		if rp.serverConfig == nil {
			config, err := certs.LoadServerCerts(rp.rootCert, rp.serverCert, rp.serverKey)
			if err != nil {
				return err
			}
			rp.serverConfig = config
		}
		// Start listening through TLS protocol
		return rp.startTLS()
	default:
		return errors.New("listen protocol not supported")
	}
}

func (rp *RProxy) serve(conn net.Conn) error {
	switch rp.backendProto {
	case "tcp":
		return rp.serveTCP(conn)
	case "tls":
		return rp.serveTLS(conn)
	default:
		return errors.New("backend protocol not supported")
	}
}

func (rp *RProxy) startTCP() error {
	// Resolve network address
	lAddr, err := net.ResolveTCPAddr("tcp", rp.listenAddr)
	if err != nil {
		return err
	}
	// Listen to TCP connections
	ln, err := net.ListenTCP("tcp", lAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	// Handle connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go func() {
			if err := rp.serve(conn); err != nil {
				log.Printf("serve error: %v", err)
			}
		}()
	}
}

func (rp *RProxy) startTLS() error {
	// Listen to TLS connections
	ln, err := tls.Listen("tcp", rp.listenAddr, rp.serverConfig)
	if err != nil {
		return err
	}
	defer ln.Close()
	// Handle connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error (%v)", err)
			continue
		}
		if err := rp.serve(conn); err != nil {
			log.Printf("serve error: %v", err)
		}
	}
}

func (rp *RProxy) serveTCP(listenConn net.Conn) error {
	// Dial to the backend server
	backendConn, err := net.DialTimeout("tcp", rp.backendAddr, 30*time.Second)
	if err != nil {
		listenConn.Close()
		return err
	}
	// Copy network traffic from the listen connection to backend connection
	go func() {
		io.Copy(backendConn, NewRPReader(listenConn))
		backendConn.Close()
		listenConn.Close()
	}()
	// Copy network traffic from the backend connection to listen connection
	io.Copy(NewRPWriteCloser(listenConn), backendConn)
	backendConn.Close()
	listenConn.Close()
	return nil
}

func (rp *RProxy) serveTLS(listenConn net.Conn) error {
	// Dial to the beckend server
	backendConn, err := tls.Dial("tcp", rp.backendAddr, rp.clientConfig)
	if err != nil {
		listenConn.Close()
		return err
	}
	// Copy network traffic from the listen connection to backend connection
	go func() {
		io.Copy(backendConn, NewRPReader(listenConn))
		backendConn.Close()
		listenConn.Close()
	}()
	// Copy network traffic from the backend connection to listen connection
	io.Copy(NewRPWriteCloser(listenConn), backendConn)
	backendConn.Close()
	listenConn.Close()
	return nil
}
