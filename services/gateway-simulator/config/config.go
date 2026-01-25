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

	TenantIDs []string `envconfig:"TENANT_IDS" default:"tenant-001,tenant-002"`

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
