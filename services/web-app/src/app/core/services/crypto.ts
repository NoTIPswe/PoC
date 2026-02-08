import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class CryptoService {
  private readonly HEX_KEY = '0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef';
  private cryptoKey: CryptoKey | null = null;
  private keyInitialized: Promise<void>;

  constructor() {
    this.keyInitialized = this.initKey();
  }

  private async initKey(): Promise<void> {
    const keyBytes = this.hexToBytes(this.HEX_KEY);
    this.cryptoKey = await window.crypto.subtle.importKey(
      'raw',
      keyBytes as any,
      { name: 'AES-GCM' },
      false,
      ['decrypt'],
    );
  }

  async decryptPayload(nonceBase64: string, ciphertextBase64: string): Promise<string> {
    await this.keyInitialized;

    try {
      const iv = this.base64ToBytes(nonceBase64);
      const data = this.base64ToBytes(ciphertextBase64);

      const decryptedBuffer = await window.crypto.subtle.decrypt(
        { name: 'AES-GCM', iv: iv as any },
        this.cryptoKey!,
        data as any,
      );

      return new TextDecoder().decode(decryptedBuffer);
    } catch (error) {
      console.error('Decryption failed:', error);
      throw error;
    }
  }

  private hexToBytes(hex: string): Uint8Array {
    const bytes = new Uint8Array(hex.length / 2);
    for (let i = 0; i < hex.length; i += 2) bytes[i / 2] = parseInt(hex.substring(i, i + 2), 16);
    return bytes;
  }

  private base64ToBytes(base64: string): Uint8Array {
    const binaryString = window.atob(base64);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) bytes[i] = binaryString.charCodeAt(i);
    return bytes;
  }
}
