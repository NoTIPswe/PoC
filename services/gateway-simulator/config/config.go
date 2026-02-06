package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NATSURL string `envconfig:"NATS_URL" default:"tls://localhost:4222"`

	TLSCACert     string `envconfig:"TLS_CA_CERT" default:"certs/ca.crt"`
	TLSClientCert string `envconfig:"TLS_CLIENT_CERT" default:"certs/client.crt"`
	TLSClientKey  string `envconfig:"TLS_CLIENT_KEY" default:"certs/client.key"`

	NumGateways       int           `envconfig:"NUM_GATEWAYS" default:"2"`
	DevicesPerGateway int           `envconfig:"DEVICES_PER_GATEWAY" default:"5"`
	TelemetryInterval time.Duration `envconfig:"TELEMETRY_INTERVAL" default:"5s"`

	TenantIDs []string `envconfig:"TENANT_IDS" default:"605e76a6-9812-4632-8418-43d99d9403d1,a66b9370-13f8-43e3-b097-f58c704f0f62"`

	EncryptionKey string `envconfig:"ENCRYPTION_KEY" default:"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"`

	LogLevel         string        `envconfig:"LOG_LEVEL" default:"info"`
	GracefulShutdown time.Duration `envconfig:"GRACEFUL_SHUTDOWN" default:"10s"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
