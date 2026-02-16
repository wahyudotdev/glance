package ca

import (
	"crypto/tls"
	"crypto/x509"
	"log"

	"github.com/elazarl/goproxy"
)

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
}
