import { Controller, DefaultValuePipe, Get, Param, ParseDatePipe, ParseIntPipe, Query, UseGuards } from "@nestjs/common";
import { TelemetryService } from "./telemetry.service";
import { TenantGuard } from "../common/guards/tenant.guard";
import { TenantId } from "../common/decorators/tenant-id.decorator";

@Controller('telemetry')
@UseGuards(TenantGuard)
export class TelemetryController {
    constructor(private readonly telemetryService: TelemetryService) {}

    @Get()
    async list(
        @TenantId() tenantId: string, 
        @Query('limit', new DefaultValuePipe(100), ParseIntPipe) limit: number, 
        @Query('offset', new DefaultValuePipe(0), ParseIntPipe) offset: number,
        @Query('from', new ParseDatePipe({ optional: true})) from?: Date,
        @Query('to', new ParseDatePipe({ optional: true})) to?: Date
    ) {
        return this.telemetryService.findByTenant(
            tenantId, 
            limit, 
            offset,
            from, 
            to
        );
    }

    @Get('count')
    async count(@TenantId() tenantId: string) {
        return { count: await this.telemetryService.countByTenant(tenantId) };
    }

    @Get('latest')
    async latest(@TenantId() tenantId: string) {
        return this.telemetryService.getLatest(tenantId);
    }

    @Get('gateways')
    async gateways(@TenantId() tenantId: string) {
        return this.telemetryService.getGateways(tenantId);
    }

    @Get('gateway/:gatewayId')
    async listByGateway(
        @TenantId() tenantId: string, 
        @Param('gatewayId') gatewayId: string, 
        @Query('limit', new DefaultValuePipe(100), ParseIntPipe) limit: number
    ) {
        return this.telemetryService.findByGateway(
            tenantId, 
            gatewayId, 
            limit
        );
    }
}