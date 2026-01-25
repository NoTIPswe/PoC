CREATE EXTENSION IF NOT EXISTS timescaledb;

-- ============================================
-- CONFIGURAZIONE
-- ============================================

CREATE TABLE IF NOT EXISTS tenants (
    id               TEXT PRIMARY KEY,
    name             TEXT NOT NULL,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    active           BOOLEAN DEFAULT TRUE
);

-- ============================================
-- TELEMETRIA
-- ============================================

CREATE TABLE IF NOT EXISTS telemetry_envelopes (
    time             TIMESTAMPTZ NOT NULL,
    tenant_id        TEXT NOT NULL,
    gateway_id       TEXT NOT NULL,
    version          INTEGER NOT NULL DEFAULT 1,
    key_id           TEXT NOT NULL,
    nonce            TEXT NOT NULL,
    ciphertext       TEXT NOT NULL
);

SELECT create_hypertable('telemetry_envelopes', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_tenant_time ON telemetry_envelopes (tenant_id, time DESC);
CREATE INDEX IF NOT EXISTS idx_gateway_time ON telemetry_envelopes (gateway_id, time DESC);
