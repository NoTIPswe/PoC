import { Column, Entity, PrimaryColumn } from "typeorm";

@Entity('telemetry_envelopes')
export class TelemetryEnvelope {
    @PrimaryColumn({ type: 'bigint' })
    id: number;

    @PrimaryColumn({ type: 'timestamptz' })
    time: Date;

    @Column({ type: 'uuid', name: 'tenant_id' })
    tenantId: string;

    @Column({ type: 'uuid', name: 'gateway_id' })
    gatewayId: string;

    @Column({ 
        type: 'bytea',
        transformer: {
            to: (value: string) => Buffer.from(value, 'base64'),
            from: (value: Buffer) => value?.toString('base64')
        }
    })
    nonce: string;

    @Column({ type: 'int' })
    version: number;

    @Column({ type: 'text', name: 'key_id' })
    keyId: string;

    @Column({ 
        type: 'bytea',
        transformer: {
            to: (value: string) => Buffer.from(value, 'base64'),
            from: (value: Buffer) => value?.toString('base64')
        }
    })
    ciphertext: string;
}