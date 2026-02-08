import { Component, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TelemetryService } from '../core/services/telemetry';
import { CryptoService } from '../core/services/crypto';
import { DecryptedTelemetry, EncryptedEnvelope } from '../core/interfaces/telemetry';
import { Subscription, timer } from 'rxjs';
import { switchMap } from 'rxjs/operators';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './dashboard.html',
  styleUrls: ['./dashboard.scss'],
})
export class DashboardComponent implements OnDestroy {
  isLoggedIn = false;
  tenantIdInput = '';
  currentTenant = '';

  telemetryData: DecryptedTelemetry[] = [];
  isLoading = false;
  errorMsg = '';

  private pollSub: Subscription | null = null;

  constructor(
    private telemetryService: TelemetryService,
    private cryptoService: CryptoService,
    private cdr: ChangeDetectorRef,
  ) {}

  login() {
    if (!this.tenantIdInput) return;

    this.currentTenant = this.tenantIdInput.trim();
    this.isLoggedIn = true;
    this.isLoading = true;
    this.errorMsg = '';

    this.startPolling();
  }

  logout() {
    this.isLoggedIn = false;
    this.currentTenant = '';
    this.tenantIdInput = '';
    this.telemetryData = [];
    this.errorMsg = '';
    this.stopPolling();
  }

  private startPolling() {
    this.pollSub = timer(0, 2000)
      .pipe(
        switchMap(() => {
          console.log('Fetching telemetry...');
          return this.telemetryService.getLatestTelemetry(this.currentTenant);
        }),
      )
      .subscribe({
        next: async (envelopes) => {
          console.log('Received envelopes:', envelopes);
          this.errorMsg = '';

          try {
            await this.processEnvelopes(envelopes);
            console.log('Decrypted data:', this.telemetryData);
          } catch (error) {
            console.error('Processing error:', error);
            this.errorMsg = 'Errore nella decifrazione dei dati';
          }

          this.isLoading = false;
          this.cdr.detectChanges();
        },
        error: (err) => {
          console.error('API Error:', err);
          this.isLoading = false;
          this.errorMsg = 'Errore di connessione al server API';
          this.cdr.detectChanges();
        },
      });
  }

  private stopPolling() {
    if (this.pollSub) {
      this.pollSub.unsubscribe();
      this.pollSub = null;
    }
  }

  ngOnDestroy(): void {
    this.stopPolling();
  }

  private async processEnvelopes(envelopes: EncryptedEnvelope[]): Promise<void> {
    if (!envelopes || envelopes.length === 0) {
      console.log('No envelopes received');
      this.telemetryData = [];
      return;
    }

    const decryptedList: DecryptedTelemetry[] = [];

    for (const env of envelopes) {
      try {
        const jsonString = await this.cryptoService.decryptPayload(env.nonce, env.ciphertext);
        const data: DecryptedTelemetry = JSON.parse(jsonString);

        decryptedList.push(data);
      } catch (e) {
        console.error('Errore decifrazione per envelope:', env, e);
      }
    }

    const existingIds = new Set(this.telemetryData.map((item) => item.message_id));

    const newItems = decryptedList.filter((item) => !existingIds.has(item.message_id));

    this.telemetryData = [...this.telemetryData, ...newItems];

    this.telemetryData.sort(
      (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime(),
    );

    if (this.telemetryData.length > 500) {
      this.telemetryData = this.telemetryData.slice(0, 500);
    }

    console.log(
      `Final telemetryData: ${this.telemetryData.length} elementi totali (${newItems.length} nuovi)`,
    );
  }

  formatMeasurements(measurement: any): string {
    if (!measurement) return 'N/A';
    return Object.entries(measurement)
      .map(([k, v]) => `${k}: ${v}`)
      .join(', ');
  }
}
