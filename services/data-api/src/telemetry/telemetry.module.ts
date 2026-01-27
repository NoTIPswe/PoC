import { Module } from "@nestjs/common";
import { TelemetryService } from "./telemetry.service";
import { TypeOrmModule } from "@nestjs/typeorm";
import { TelemetryEnvelope } from "./entities/telemetry-envelope.entity";
import { TelemetryController } from "./telemetry.controller";

@Module({
    imports: [TypeOrmModule.forFeature([TelemetryEnvelope])], 
    controllers: [TelemetryController],
    providers: [TelemetryService],
    exports: []
})
export class TelemetryModule {}