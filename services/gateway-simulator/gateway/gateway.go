package gateway

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gateway-simulator/crypto"
	"gateway-simulator/device"
	"gateway-simulator/telemetry"
	"gateway-simulator/transport"
)

type Gateway struct {
	ID        string
	TenantID  string
	devices   []*device.Device
	transport *transport.Client
	encryptor *crypto.Encryptor
}

func NewGateway(id, tenantID string, numDevices int, client *transport.Client, encryptor *crypto.Encryptor) *Gateway {
	gw := &Gateway{
		ID:        id,
		TenantID:  tenantID,
		devices:   make([]*device.Device, 0, numDevices),
		transport: client,
		encryptor: encryptor,
	}

	// Crea dispositivi con tipi diversi in round-robin
	deviceTypes := device.GetDeviceTypes()
	for i := 0; i < numDevices; i++ {
		devType := deviceTypes[i%len(deviceTypes)]
		dev := device.NewDevice(devType, i+1, id, tenantID)
		gw.devices = append(gw.devices, dev)
	}

	return gw
}

func (g *Gateway) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("[Gateway %s] Started with %d devices (tenant: %s)", g.ID, len(g.devices), g.TenantID)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Gateway %s] Shutting down", g.ID)
			return
		case <-ticker.C:
			g.publishTelemetry()
		}
	}
}

func (g *Gateway) publishTelemetry() {
	for _, dev := range g.devices {
		payload := dev.GenerateTelemetry()

		encrypted, err := g.encryptPayload(payload)
		if err != nil {
			log.Printf("[Gateway %s] Encryption error: %v", g.ID, err)
			continue
		}

		subject := fmt.Sprintf("telemetry.%s.%s.%s", g.TenantID, g.ID, payload.DeviceType)

		if err := g.transport.Publish(subject, encrypted); err != nil {
			log.Printf("[Gateway %s] Publish error: %v", g.ID, err)
			continue
		}
	}
}

func (g *Gateway) encryptPayload(payload telemetry.Payload) ([]byte, error) {
	plaintext, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	nonce, ciphertext, err := g.encryptor.Encrypt(plaintext)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	envelope := telemetry.EncryptedEnvelope{
		Version:    1,
		TenantID:   payload.TenantID,
		GatewayID:  payload.GatewayID,
		KeyID:      g.encryptor.KeyID(),
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		Timestamp:  payload.Timestamp,
	}

	return json.Marshal(envelope)
}

func (g *Gateway) DeviceCount() int {
	return len(g.devices)
}
