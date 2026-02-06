package messaging

import (
	"log"

	"github.com/nats-io/nats.go"
)

type NATSClient struct {
	Conn *nats.Conn
	JS   nats.JetStreamContext
}

func InitNATS(url string, tlsOpt nats.Option) *NATSClient {
	nc, err := nats.Connect(url,
		tlsOpt,
		nats.Name("data-consumer"),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			if err != nil {
				log.Printf("[NATS] Disconnected %v", err)
			}
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Printf("[NATS] Reconnected to %s", c.ConnectedUrl())
		}),
	)

	if err != nil {
		log.Fatalf("NATS connection failed %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Jetstream init failed %v", err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "TELEMETRY",
		Subjects:  []string{"telemetry.>"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
	})
	if err != nil {
		log.Fatalf("Stream creation failed: %v", err)
	}

	log.Println("[NATS] Connected, Jetstream stream TELEMETRY ready")
	return &NATSClient{Conn: nc, JS: js}
}
