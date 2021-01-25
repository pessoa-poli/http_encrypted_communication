package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var ()

func main() {
	clientCertFile := "/home/luado/lp/http_encrypted_communication/certs/out/client.crt"
	clientKeyFile := "/home/luado/lp/http_encrypted_communication/certs/out/client.key"
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", clientCertFile, clientKeyFile)
	}
	caCertFile := "/home/luado/lp/http_encrypted_communication/certs/out/rnpCA.crt"
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Error opening cert file %s, Error: %s", caCertFile, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}

	client := http.Client{Transport: t, Timeout: 15 * time.Second}
	resp, err := client.Get("https://localhost:9000")
	if err != nil {
		panic("Failed to do request: " + err.Error())
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("\nServer replied with: %s\n", respBytes)
}
