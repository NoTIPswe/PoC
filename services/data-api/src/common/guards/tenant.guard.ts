import {
  BadRequestException,
  CanActivate,
  ExecutionContext,
  Injectable,
} from '@nestjs/common';

@Injectable()
export class TenantGuard implements CanActivate {
  canActivate(context: ExecutionContext): boolean {
    const request = context.switchToHttp().getRequest<Request>();
    const tenantId = request.headers['x-tenant-id'];

    if (!tenantId || Array.isArray(tenantId)) {
      throw new BadRequestException('X-Tenant-Id header is required');
    }

    return true;
  }
}
