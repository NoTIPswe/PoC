import { Component, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TelemetryService } from '../core/services/telemetry'; 
import { CryptoService } from '../core/services/crypto';       
import { DecryptedTelemetry, EncryptedEnvelope } from '../core/interfaces/telemetry'; 
import { Subscription, switchMap, timer } from 'rxjs'; // Usa 'timer' invece di 'interval'

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
    this.isLoading = true; // Mostra il caricamento iniziale
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
    // timer(0, 2000) -> Parte subito (0ms) e poi ripete ogni 2000ms
    // Risolve il problema dell'attesa iniziale di 2 secondi
    this.pollSub = timer(0, 2000) 
      .pipe(
        switchMap(() => this.telemetryService.getLatestTelemetry(this.currentTenant))
      )
      .subscribe({
        next: async (envelopes) => { // Nota 'async' qui
          this.errorMsg = '';
          // Aspettiamo che la decifrazione sia finita PRIMA di spegnere il caricamento
          await this.processEnvelopes(envelopes);
          this.isLoading = false; 
        },
        error: (err) => {
           console.error('API Error:', err);
           this.isLoading = false;
           this.errorMsg = 'Errore di connessione al server API';
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
    // Controllo di sicurezza: se envelopes è null/undefined, fermati
    if (!envelopes) return;

    const decryptedList: DecryptedTelemetry[] = [];
    
    for (const env of envelopes) {
      try {
        const jsonString = await this.cryptoService.decryptPayload(env.nonce, env.ciphertext);
        const data = JSON.parse(jsonString);

        // Mappatura corretta per l'HTML
        data.time = env.time; 
        data.tenant_id = env.tenantId;

        decryptedList.push(data);
      } catch (e) {
        console.error('Errore decifrazione:', e);
      }
    }

    // Aggiorna la tabella e ordina
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