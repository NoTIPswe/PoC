import { Component, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TelemetryService } from '../core/services/telemetry'; 
import { CryptoService } from '../core/services/crypto';       
import { DecryptedTelemetry, EncryptedEnvelope } from '../core/interfaces/telemetry'; 
import { interval, Subscription, switchMap } from 'rxjs';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './dashboard.html',
  styleUrls: ['./dashboard.scss']
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
    private cryptoService: CryptoService
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

    this.pollSub = interval(2000)
      .pipe(
        switchMap(() => this.telemetryService.getLatestTelemetry(this.currentTenant))
      )
      .subscribe({
        next: (envelopes) => {
          this.processEnvelopes(envelopes);
          this.isLoading = false;
          this.errorMsg = ''; 
        },
        error: (err) => {
           console.error('API Error:', err);
           this.isLoading = false;
           this.errorMsg = 'Errore di connessione al server API (verificare se Docker è attivo)';
        }
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


  private async processEnvelopes(envelopes: EncryptedEnvelope[]) {
    const decryptedList: DecryptedTelemetry[] = [];
    
    for (const env of envelopes) {
      try {
        const jsonString = await this.cryptoService.decryptPayload(env.nonce, env.ciphertext);
        
        const data = JSON.parse(jsonString);


        data.timestamp = env.time; 
        data.tenant_id = env.tenantId;

        decryptedList.push(data);

      } catch (e) {
        console.error('Errore decifrazione per envelope:', env, e);
      }
    }

    this.telemetryData = decryptedList.sort((a, b) => 
      new Date(b.time).getTime() - new Date(a.time).getTime()
    );
  }

  formatMeasurements(measurements: any): string {
    if (!measurements) return '';
    return Object.entries(measurements)
      .map(([k, v]) => `${k}: ${v}`)
      .join(', ');
  }
}