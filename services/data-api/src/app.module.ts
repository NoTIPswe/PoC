import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { TypeOrmModule } from '@nestjs/typeorm';
import { TelemetryModule } from './telemetry/telemetry.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }), 
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: (configService: ConfigService) => ({
        type: 'postgres',
        host: configService.get('DB_HOST', 'localhost'),
        port: configService.get('DB_PORT', 5432),
        username: configService.get('DB_USERNAME', 'poc'),
        password: configService.get('DB_PASSWORD', 'poc_password'),
        database: configService.get('DB_DATABASE', 'measures'),
        entities: [__dirname + '/**/*.entity{.ts,.js}'],
        synchronize: false    // temporarily
      }),
    }),
    TelemetryModule
  ],
})
export class AppModule {}
