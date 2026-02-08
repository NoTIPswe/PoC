import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { TelemetryEnvelope } from './entities/telemetry-envelope.entity';
import { Between, LessThanOrEqual, MoreThanOrEqual, Repository } from 'typeorm';

@Injectable()
export class TelemetryService {
  constructor(
    @InjectRepository(TelemetryEnvelope)
    private readonly telemetryRepository: Repository<TelemetryEnvelope>,
  ) {}

  async findByTenant(
    tenantId: string,
    limit: number = 100,
    offset: number = 0, // only if needed
    from?: Date,
    to?: Date,
  ): Promise<TelemetryEnvelope[]> {
    const where: any = { tenantId };
    if (from && to) {
      where.time = Between(from, to);
    } else if (from) {
      where.time = MoreThanOrEqual(from);
    } else if (to) {
      where.time = LessThanOrEqual(to);
    }

    return this.telemetryRepository.find({
      where,
      order: { time: 'DESC' },
      take: limit,
      skip: offset,
    });
  }

  async findByGateway(
    tenantId: string,
    gatewayId: string,
    limit: number = 100,
  ): Promise<TelemetryEnvelope[]> {
    return this.telemetryRepository.find({
      where: { tenantId, gatewayId },
      order: { time: 'DESC' },
      take: limit,
    });
  }

  async countByTenant(tenantId: string): Promise<number> {
    return this.telemetryRepository.count({
      where: { tenantId },
    });
  }

  async getLatest(tenantId: string): Promise<TelemetryEnvelope | null> {
    return this.telemetryRepository.findOne({
      where: { tenantId },
      order: { time: 'DESC' },
    });
  }

  async getGateways(tenantId: string): Promise<string[]> {
    const results = await this.telemetryRepository
      .createQueryBuilder('t')
      .select('DISTINCT t.gateway_id', 'gatewayId')
      .where('t.tenant_id = :tenantId', { tenantId })
      .getRawMany();
    return results.map((r) => r.gatewayId);
  }
}
