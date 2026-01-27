package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	nc := messaging.InitNATS(cfg.NATSURL, nats.Secure(tlsCfg))
	log.Println("✅ Connesso a NATS")
	defer nc.Close()

	dbPool := storage.InitDatabase(ctx, cfg.DBURL)
	log.Println("✅ Connesso a TimescaleDB")
	defer dbPool.Close()

	_, err = nc.Subscribe("telemetry.>", func(m *nats.Msg) {
		var t storage.Telemetry
		if err := json.Unmarshal(m.Data, &t); err != nil {
			log.Printf("JSON Error: %v", err)
			return
		}

		if err := storage.SaveTelemetry(ctx, dbPool, t); err != nil {
			log.Printf("DB Error: %v", err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Data Consumer avviato correttamente")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
