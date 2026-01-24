package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type TLSConfig struct {
	CACert     string
	ClientCert string
	ClientKey  string
}

type Client struct {
	client paho.Client
	qos    byte
}

func NewClient(broker, clientID string, qos byte, tlsCfg *TLSConfig) (*Client, error) {
	opts := paho.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetConnectionLostHandler(func(c paho.Client, err error) {
			log.Printf("[MQTT] Connection lost: %v", err)
		}).
		SetOnConnectHandler(func(c paho.Client) {
			log.Printf("[MQTT] Connected: %s", clientID)
		})

	if tlsCfg != nil {
		tlsConfig, err := loadTLSConfig(tlsCfg)
		if err != nil {
			return nil, fmt.Errorf("load TLS config: %w", err)
		}
		opts.SetTLSConfig(tlsConfig)
	}

	client := paho.NewClient(opts)

	token := client.Connect()
	if !token.WaitTimeout(10 * time.Second) {
		return nil, fmt.Errorf("connection timeout")
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	return &Client{client: client, qos: qos}, nil
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
	token := c.client.Publish(topic, c.qos, false, payload)
	token.Wait()
	return token.Error()
}

func (c *Client) Disconnect() {
	c.client.Disconnect(1000)
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}
