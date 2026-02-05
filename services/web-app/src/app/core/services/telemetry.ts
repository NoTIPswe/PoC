import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { EncryptedEnvelope } from '../interfaces/telemetry';

@Injectable({
  providedIn: 'root',
})
export class TelemetryService {
  // --- MODIFICA QUI ---
  // Togliamo "http://localhost:3000". Usiamo il percorso relativo.
  // Il browser chiamerà http://localhost:4200/api/v1/telemetry
  // e Nginx lo girerà al backend.
  private apiUrl = '/api/v1/telemetry'; 

  constructor(private http: HttpClient) {}

  getLatestTelemetry(tenantId: string): Observable<EncryptedEnvelope[]> {
    // Mantieni l'header che abbiamo aggiunto prima
    const headers = new HttpHeaders().set('x-tenant-id', tenantId);
    return this.http.get<EncryptedEnvelope[]>(this.apiUrl, { headers });
  }
}