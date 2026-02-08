export interface EncryptedEnvelope {
  time: string;
  tenantId: string;
  gatewayId: string;
  version: number;
  keyId: string;
  nonce: string;
  ciphertext: string;
}

export interface DecryptedTelemetry {
  message_id: string;
  gateway_id: string;
  device_id: string;
  tenant_id: string;
  device_type: string;
  timestamp: string;
  measurement: { [key: string]: number | boolean };
}
