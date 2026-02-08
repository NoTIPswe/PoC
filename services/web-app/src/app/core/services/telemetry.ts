import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { EncryptedEnvelope } from '../interfaces/telemetry';

@Injectable({
  providedIn: 'root',
})
export class TelemetryService {
  private apiUrl = '/api/v1/telemetry';

  constructor(private http: HttpClient) {}

  getLatestTelemetry(tenantId: string): Observable<EncryptedEnvelope[]> {
    const headers = new HttpHeaders().set('x-tenant-id', tenantId);
    return this.http.get<EncryptedEnvelope[]>(this.apiUrl, { headers });
  }
}
