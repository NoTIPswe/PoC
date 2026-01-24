package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"gateway-simulator/config"
	"gateway-simulator/crypto"
	"gateway-simulator/gateway"
)

func main() {
	log.Println("Gateway Simulator starting...")

	// Carica configurazione
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Inizializza encryption
	encryptor, err := crypto.NewEncryptor(cfg.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize encryption: %v", err)
	}

	// Context per graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Avvia simulatore
	sim := gateway.NewSimulator(cfg, encryptor)
	if err := sim.Start(ctx); err != nil {
		log.Fatalf("Failed to start simulator: %v", err)
	}

	// Attendi segnale di shutdown
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Graceful shutdown
	sim.Shutdown(cfg.GracefulShutdown)
}
