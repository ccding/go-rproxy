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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

const (
	appName        = "testapp"
	rootType       = "root"
	serverType     = "server"
	clientType     = "client"
	rootFileName   = appName + "-root"
	serverFileName = appName + "-server"
	clientFileName = appName + "-client-"
	keySize        = 2048
	validDuration  = 365 * 24 * time.Hour
)

var (
	certTypes = map[string]bool{
		rootType:   true,
		serverType: true,
		clientType: true,
	}
	serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)
)

func main() {
	var certType = flag.String("type", "", "certificate type: root, server, or client")
	var clientID = flag.String("cid", "", "client id: required when type==client")
	flag.Parse()
	// Check flags
	if _, ok := certTypes[*certType]; !ok {
		log.Fatalf("bad certificate type")
	}
	if *certType == clientType && len(*clientID) == 0 {
		log.Fatalf("bad client id (required)")
	}
	priv, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
	}
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}
	// Valid time
	notBefore := time.Now()
	notAfter := notBefore.Add(validDuration)
	// Set key config
	var (
		extKeyUsage []x509.ExtKeyUsage
		isCA        bool
		commonName  string
		keyUsage    x509.KeyUsage
	)
	if *certType == rootType {
		extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
		isCA = true
		commonName = rootFileName
		keyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	} else if *certType == serverType {
		extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		isCA = false
		commonName = serverFileName
		keyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	} else if *certType == clientType {
		extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
		isCA = false
		commonName = clientFileName + *clientID
		keyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	}

	var (
		rootCert *x509.Certificate
		rootKey  interface{}
	)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: commonName},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     keyUsage,
		ExtKeyUsage:  extKeyUsage,
		IsCA:         isCA,
		BasicConstraintsValid: true,
	}

	if *certType == "root" {
		// Root is a self-signed certificate
		rootCert = &template
		rootKey = priv
	} else {
		// Load root certificate/key to sign client or server certificate
		log.Print("loading root_cert.pem and root_key.pem")

		rootCertData, err := ioutil.ReadFile("root_cert.pem")
		if err != nil {
			log.Fatalf("failed to read root certificate: %s", err)
		}
		rootCertBlock, _ := pem.Decode(rootCertData)
		if rootCertBlock == nil {
			log.Fatalf("failed to decode root certificate pem")
		}
		rootCert, err = x509.ParseCertificate(rootCertBlock.Bytes)
		if err != nil {
			log.Fatalf("failed to parse root certificate: %s", err)
		}

		rootKeyData, err := ioutil.ReadFile("root_key.pem")
		if err != nil {
			log.Fatalf("failed to read root private key: %s", err)
		}
		rootKeyBlock, _ := pem.Decode(rootKeyData)
		if rootKeyBlock == nil {
			log.Fatalf("failed to decode root private key pem")
		}
		rootKey, err = x509.ParsePKCS1PrivateKey(rootKeyBlock.Bytes)
		if err != nil {
			log.Fatalf("failed to parse root private key: %s", err)
		}
	}
	// Create certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, rootCert, &priv.PublicKey, rootKey)
	if err != nil {
		log.Fatalf("failed to create certificate: %s", err)
	}
	// Get file names
	var certfn, keyfn string
	if *certType == rootType {
		certfn = "root_cert.pem"
		keyfn = "root_key.pem"
	} else if *certType == clientType {
		certfn = "client_" + *clientID + "_cert.pem"
		keyfn = "client_" + *clientID + "_key.pem"
	} else if *certType == serverType {
		certfn = "server_cert.pem"
		keyfn = "server_key.pem"
	}
	// Write to files
	certOut, err := os.Create(certfn)
	if err != nil {
		log.Fatalf("failed to open %s for writing: %s", certfn, err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certOut.Close()
	log.Println("certificate:", certfn)

	keyOut, err := os.OpenFile(keyfn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("failed to open %s for writing:", keyfn, err)
		return
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Println("private key:", keyfn)

	log.Println("common name:", commonName)
}
