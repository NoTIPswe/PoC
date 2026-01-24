package gateway

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"gateway-simulator/config"
	"gateway-simulator/crypto"
	"gateway-simulator/mqtt"
)

type Simulator struct {
	config    *config.Config
	encryptor *crypto.Encryptor
	gateways  []*Gateway
	clients   []*mqtt.Client
	wg        sync.WaitGroup
}

func NewSimulator(cfg *config.Config, encryptor *crypto.Encryptor) *Simulator {
	return &Simulator{
		config:    cfg,
		encryptor: encryptor,
		gateways:  make([]*Gateway, 0, cfg.NumGateways),
		clients:   make([]*mqtt.Client, 0, cfg.NumGateways),
	}
}

func (s *Simulator) Start(ctx context.Context) error {
	log.Printf("[Simulator] Starting %d gateways, %d devices each", s.config.NumGateways, s.config.DevicesPerGateway)

	for i := 0; i < s.config.NumGateways; i++ {
		gatewayID := fmt.Sprintf("gw-%03d", i+1)
		tenantID := s.config.TenantIDs[i%len(s.config.TenantIDs)]
		clientID := fmt.Sprintf("%s-%s", s.config.MQTTClientPrefix, gatewayID)

		tlsCfg := &mqtt.TLSConfig{
			CACert:     s.config.TLSCACert,
			ClientCert: s.config.TLSClientCert,
			ClientKey:  s.config.TLSClientKey,
		}
		mqttClient, err := mqtt.NewClient(s.config.MQTTBroker, clientID, s.config.MQTTQoS, tlsCfg)
		if err != nil {
			return fmt.Errorf("create MQTT client for %s: %w", gatewayID, err)
		}
		s.clients = append(s.clients, mqttClient)

		gw := NewGateway(gatewayID, tenantID, s.config.DevicesPerGateway, mqttClient, s.encryptor)
		s.gateways = append(s.gateways, gw)

		s.wg.Add(1)
		go func(g *Gateway) {
			defer s.wg.Done()
			g.Run(ctx, s.config.TelemetryInterval)
		}(gw)
	}

	totalDevices := s.config.NumGateways * s.config.DevicesPerGateway
	log.Printf("[Simulator] All gateways started. Total devices: %d", totalDevices)

	return nil
}

func (s *Simulator) Shutdown(timeout time.Duration) {
	log.Println("[Simulator] Shutting down...")

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[Simulator] All gateways stopped")
	case <-time.After(timeout):
		log.Println("[Simulator] Shutdown timeout, forcing exit")
	}

	for _, client := range s.clients {
		client.Disconnect()
	}

	log.Println("[Simulator] Shutdown complete")
}

func (s *Simulator) Stats() (gateways, devices int) {
	gateways = len(s.gateways)
	for _, gw := range s.gateways {
		devices += gw.DeviceCount()
	}
	return
}
