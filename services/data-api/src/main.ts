import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { VersioningType } from '@nestjs/common';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  app.enableCors({
    origin: '*',
    methods: 'GET,HEAD',
    allowedHeaders: 'Content-Type, Accept, X-Tenant-Id',
  });

  app.setGlobalPrefix('api'); 
  app.enableVersioning({
    type: VersioningType.URI,
    defaultVersion: '1',
  })

  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
