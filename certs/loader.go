package certs

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

func LoadClientCerts(rootCert, clientCert, clientKey, serverName string) (*tls.Config, error) {
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

func LoadServerCerts(rootCert, serverCert, serverKey string) (*tls.Config, error) {
	// Load root certificate to verify client certificate
	rootPEM, err := ioutil.ReadFile(rootCert)
	if err != nil {
		return nil, err
	}
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM([]byte(rootPEM)); !ok {
		return nil, errors.New("failed to parse root certificate")
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
