import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { EncryptedEnvelope } from '../interfaces/telemetry';

@Injectable({
  providedIn: 'root',
})
export class TelemetryService {
  private apiUrl = 'http://localhost:3000/api/telemetry';
  constructor(private http: HttpClient) {}
  getLatestTelemetry(tenantId: string): Observable<EncryptedEnvelope[]> {
    const params = new HttpParams().set('tenantId', tenantId);
    return this.http.get<EncryptedEnvelope[]>(this.apiUrl, { params });
  }
}