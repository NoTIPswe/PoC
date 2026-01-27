package storage

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Telemetry struct {
	Time       time.Time `json:"timestamp"`
	TenantID   string    `json:"tenant_id"`
	GatewayID  string    `json:"gateway_id"`
	Version    int       `json:"version"`
	KeyID      string    `json:"key_id"`
	Nonce      string    `json:"nonce"`
	Ciphertext string    `json:"ciphertext"`
}

func InitDatabase(ctx context.Context, url string) *pgxpool.Pool {
	dbPool, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatalf("Connessione a DB fallita: %v", err)
	}
	return dbPool
}

func SaveTelemetry(ctx context.Context, db *pgxpool.Pool, t Telemetry) error {
	query := `INSERT INTO telemetry_envelopes (
			time, tenant_id, gateway_id, version, key_id, nonce, ciphertext
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(ctx, query,
		t.Time,
		t.TenantID,
		t.GatewayID,
		t.Version,
		t.KeyID,
		t.Nonce,
		t.Ciphertext,
	)
	return err
}
