package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NATSURL string `envconfig:"NATS_URL" default:"tls://nats:4222"`
	DBURL   string `enconfig:"DB_URL" default:"postgres://poc:poc_password@timescaledb:5432/measures"`

	TLSCACert     string `envconfig:"TLS_CA_CERT" default:"certs/ca.crt"`
	TLSClientCert string `envconfig:"TLS_CLIENT_CERT" default:"certs/client.crt"`
	TLSClientKey  string `envconfig:"TLS_CLIENT_KEY" default:"certs/client.key"`

	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
