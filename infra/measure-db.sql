CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================
-- CONFIGURAZIONE
-- ============================================

CREATE TABLE IF NOT EXISTS tenants (
    id               UUID PRIMARY KEY DEFAULT uuidv4(),
    name             TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS gateway (
    id              UUID PRIMARY KEY DEFAULT uuidv4()
);

CREATE TABLE IF NOT EXISTS tenants_gateways (
    id_tenant       UUID NOT NULL,
    id_gateway      UUID NOT NULL,

    FOREIGN KEY (id_tenant) REFERENCES tenants(id),
    FOREIGN KEY (id_gateway) REFERENCES gateway(id),
    PRIMARY KEY (id_tenant, id_gateway)
);

-- ============================================
-- TELEMETRIA
-- ============================================

CREATE TABLE IF NOT EXISTS telemetry_envelopes (
    time             TIMESTAMPTZ NOT NULL,
    tenant_id        UUID NOT NULL DEFAULT uuidv4(),
    gateway_id       UUID NOT NULL DEFAULT uuidv4(),
    version          INTEGER NOT NULL DEFAULT 1,
    key_id           TEXT NOT NULL,
    nonce            TEXT NOT NULL,
    ciphertext       TEXT NOT NULL

    PRIMARY KEY (time, tenant_id, gateway_id, nonce)
);

SELECT create_hypertable('telemetry_envelopes', 'time', if_not_exists => TRUE, migrate_data => TRUE);

CREATE INDEX IF NOT EXISTS idx_tenant_time ON telemetry_envelopes (tenant_id, time DESC);
CREATE INDEX IF NOT EXISTS idx_gateway_time ON telemetry_envelopes (gateway_id, time DESC);
