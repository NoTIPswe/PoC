package messaging

import (
	"log"

	"github.com/nats-io/nats.go"
)

func InitNATS(url string, tlsConfig nats.Option) *nats.Conn {
	nc, err := nats.Connect(url, tlsConfig)
	if err != nil {
		log.Fatalf("Connessione NATS fallita: %v", err)
	}
	return nc
}
