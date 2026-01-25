package transport

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

type TLSConfig struct {
	CACert     string
	ClientCert string
	ClientKey  string
}

type Client struct {
	conn *nats.Conn
}

func NewClient(natsURL string, tlsConfig *TLSConfig) (*Client, error) {
	opts := []nats.Option{
		nats.Name("gateway-simulator"),
		nats.MaxReconnects(-1), // keep trying reconnection
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			if err != nil {
				log.Printf("[NATS] Disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Printf("[NATS] Recconected to %s", c.ConnectedUrl())
		}),
	}

	if tlsConfig != nil {
		tlsCfg, err := loadTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("load TLS config: %w", err)
		}
		opts = append(opts, nats.Secure(tlsCfg))
	}

	conn, err := nats.Connect(natsURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	log.Printf("[NATS] Connected to %s", conn.ConnectedUrl())
	return &Client{conn: conn}, nil
}

func loadTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	caCert, err := os.ReadFile(cfg.CACert)
	if err != nil {
		return nil, fmt.Errorf("read CA cert: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA cert")
	}

	clientCert, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
	if err != nil {
		return nil, fmt.Errorf("load client cert: %w", err)
	}

	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{clientCert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func (c *Client) Publish(topic string, payload []byte) error {
	return c.conn.Publish(topic, payload)
}

func (c *Client) Disconnect() {
	if err := c.conn.Drain(); err != nil {
		log.Printf("[NATS] Drain error: %v", err)
	}
}

func (c *Client) IsConnected() bool {
	return c.conn.IsConnected()
}
