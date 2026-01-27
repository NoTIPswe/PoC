package config

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

func LoadTLSConfig(cfg *Config) *tls.Config {
	cert, err := tls.LoadX509KeyPair(cfg.TLSClientCert, cfg.TLSClientKey)
	if err != nil {
		log.Fatalf("Errore certificati client: %v", err)
	}

	caPool := x509.NewCertPool()
	caData, err := os.ReadFile(cfg.TLSCACert)
	if err != nil {
		log.Fatalf("Errore lettura CA: %v", err)
	}
	caPool.AppendCertsFromPEM(caData)

	return &tls.Config{
		RootCAs:      caPool,
		Certificates: []tls.Certificate{cert},
	}
}
