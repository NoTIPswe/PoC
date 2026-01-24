package telemetry

import "time"

// Real payload
type Payload struct {
	MessageID   string                 `json:"message_id"`
	GatewayID   string                 `json:"gateway_id"`
	DeviceID    string                 `json:"device_id"`
	TenantID    string                 `json:"tenant_id"`
	DeviceType  string                 `json:"device_type"`
	Timestamp   time.Time              `json:"timestamp"`
	Measurements map[string]interface{} `json:"measurement"`
}

type EncryptedEnvelope struct {
	Version    int       `json:"version"`
	TenantID   string    `json:"tenant_id"`
	GatewayID  string    `json:"gateway_id"`
	KeyID      string    `json:"key_id"`
	Nonce      string    `json:"nonce"`
	Ciphertext string    `json:"ciphertext"`
	Timestamp  time.Time `json:"timestamp"`
}
