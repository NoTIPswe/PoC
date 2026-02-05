package storage

import (
	"context"
	"fmt"
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
		log.Fatalf("DB connection failed: %v", err)
	}
	return dbPool
}

type BatchWriter struct {
	db        *pgxpool.Pool
	ch        chan Telemetry
	batchSize int
	interval  time.Duration
}

func NewBatchWriter(db *pgxpool.Pool, batchSize int, interval time.Duration) *BatchWriter {
	return &BatchWriter{
		db:        db,
		ch:        make(chan Telemetry, batchSize*2),
		batchSize: batchSize,
		interval:  interval,
	}
}

func (bw *BatchWriter) Enqueue(t Telemetry) {
	bw.ch <- t
}

func (bw *BatchWriter) Run(ctx context.Context) {
	buf := make([]Telemetry, 0, bw.batchSize)
	ticker := time.NewTicker(bw.interval)
	defer ticker.Stop()

	for {
		select {
		case t := <-bw.ch:
			buf = append(buf, t)
			if len(buf) >= bw.batchSize {
				bw.flush(ctx, &buf)
			}
		case <-ticker.C:
			if len(buf) > 0 {
				bw.flush(ctx, &buf)
			}
		case <-ctx.Done():
			if len(buf) > 0 {
				bw.flush(ctx, &buf)
			}
			return
		}
	}
}

func (bw *BatchWriter) flush(ctx context.Context, buf *[]Telemetry) {
	query := "INSERT INTO telemetry_envelopes (time, tenant_id, gateway_id, version, key_id, nonce, ciphertext) VALUES "
	args := make([]interface{}, 0, len(*buf)*7)

	for i, t := range *buf {
		offset := i * 7
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7)
		args = append(args, t.Time, t.TenantID, t.GatewayID, t.Version, t.KeyID, t.Nonce, t.Ciphertext)
	}

	_, err := bw.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("[DB] Batch insert error (%d records): %v", len(*buf), err)
		return
	}

	log.Printf("[DB] Flushed %d records", len(*buf))
	*buf = (*buf)[:0]
}
