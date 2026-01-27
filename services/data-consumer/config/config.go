package config

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

func LoadTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair(os.Getenv("TLS_CLIENT_CERT"), os.Getenv("TLS_CLIENT_KEY"))
	if err != nil {
		log.Fatalf("Errore certificati client: %v", err)
	}

	caPool := x509.NewCertPool()
	caData, err := os.ReadFile(os.Getenv("TLS_CA_CERT"))
	if err != nil {
		log.Fatalf("Errore lettura CA: %v", err)
	}
	caPool.AppendCertsFromPEM(caData)

	return &tls.Config{
		RootCAs:      caPool,
		Certificates: []tls.Certificate{cert},
	}
}
