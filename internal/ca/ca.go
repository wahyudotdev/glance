package ca

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/elazarl/goproxy"
)

var CAPath string

func SetupCA() {
	// For now, we use the default goproxy CA.
	// In a production-like app, we would generate a custom one and save it to disk.
	caCert, err := tls.X509KeyPair([]byte(goproxy.CA_CERT), []byte(goproxy.CA_KEY))
	if err != nil {
		log.Fatalf("Failed to load default CA: %v", err)
	}
	if caCert.Leaf, err = x509.ParseCertificate(caCert.Certificate[0]); err != nil {
		log.Fatalf("Failed to parse default CA certificate: %v", err)
	}

	goproxy.GoproxyCa = caCert
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&caCert)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&caCert)}

	// Save CA to a temporary file
	tmpDir := os.TempDir()
	CAPath = filepath.Join(tmpDir, "agent-proxy-ca.crt")
	err = ioutil.WriteFile(CAPath, []byte(goproxy.CA_CERT), 0644)
	if err != nil {
		log.Printf("Warning: Failed to save CA certificate to %s: %v", CAPath, err)
	} else {
		log.Printf("CA certificate saved to %s", CAPath)
	}
}
