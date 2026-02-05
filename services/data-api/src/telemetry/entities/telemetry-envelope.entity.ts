import { Column, Entity, PrimaryColumn } from "typeorm";

@Entity('telemetry_envelopes')
export class TelemetryEnvelope {
    @PrimaryColumn({ type: 'timestamptz'})
    time: Date;

    @Column({ type: 'uuid', name: 'tenant_id'})
    tenantId: string;

    @Column({ type: 'uuid', name: 'gateway_id'})
    gatewayId: string;

    @PrimaryColumn({ type: 'text' })
    nonce: string;

    @Column({ type: 'int'})
    version: number;

    @Column({ type: 'text', name: 'key_id'})
    keyId: string;

    @Column({ type: 'text'})
    ciphertext: string;
}