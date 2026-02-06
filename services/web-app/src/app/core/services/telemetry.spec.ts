import { TestBed } from '@angular/core/testing';

import { TelemetryService } from './telemetry';

describe('TelemetryService', () => {
  let service: TelemetryService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TelemetryService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
