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

var (
	port    = "9000"
	host    = "localhost"
	caCert  = "/home/luado/lp/http_encrypted_communication/certs/out/rnpCA.crt"
	certOpt = tls.RequireAndVerifyClientCert
	server  = &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Minute, // 5 min to allow for delays when 'curl' on OSx prompts for username/password
		WriteTimeout: 10 * time.Second,
		TLSConfig:    getTLSConfig(host, caCert, tls.ClientAuthType(certOpt)),
	}
)

func getTLSConfig(host, caCertFile string, certOpt tls.ClientAuthType) *tls.Config {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool
	if certOpt > tls.RequestClientCert {
		caCert, err = ioutil.ReadFile(caCertFile)
		if err != nil {
			log.Fatal("Error opening cert file", caCertFile, ", error ", err)
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	return &tls.Config{
		ServerName: host,
		ClientAuth: certOpt,
		ClientCAs:  caCertPool,
		MinVersion: tls.VersionTLS12, // TLS versions below 1.2 are considered insecure - see https://www.rfc-editor.org/rfc/rfc7525.txt for details
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			body = []byte(fmt.Sprintf("error reading request body: %s", err))
		}
		resp := fmt.Sprintf("Hello, %s from Simple Server!", body)
		w.Write([]byte(resp))
	})
	serverCert := "/home/luado/lp/http_encrypted_communication/certs/out/localhost.crt"
	srvKey := "/home/luado/lp/http_encrypted_communication/certs/out/localhost.key"
	if err := server.ListenAndServeTLS(serverCert, srvKey); err != nil {
		log.Fatal(err)
	}
}
