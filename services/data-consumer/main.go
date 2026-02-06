package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"data-consumer/config"
	"data-consumer/messaging"
	"data-consumer/storage"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	tlsCfg := config.LoadTLSConfig(cfg)
	natsClient := messaging.InitNATS(cfg.NATSURL, nats.Secure(tlsCfg))
	log.Println("Connesso a NATS")
	defer natsClient.Conn.Drain() // drain all the messages

	dbPool := storage.InitDatabase(ctx, cfg.DBURL)
	log.Println("Connesso a TimescaleDB")
	defer dbPool.Close()

	writer := storage.NewBatchWriter(dbPool, 50, 2*time.Second)
	go writer.Run(ctx)


	_, err = natsClient.JS.Subscribe("telemetry.>", func(m *nats.Msg) {
		var t storage.Telemetry
		if err := json.Unmarshal(m.Data, &t); err != nil {
			log.Printf("JSON Error: %v", err)
			m.Ack()
			return
		}

		writer.Enqueue(t)
		m.Ack()

	}, nats.Durable("data-consumer"), nats.DeliverAll())

	log.Println("Data Consumer avviato correttamente (Jetstream)")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
