# Stage 1 - Base
FROM node:22-alpine AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app

# Stage 2 - Build
FROM base AS builder

COPY pnpm-lock.yaml pnpm-workspace.yaml ./
COPY services/data-api ./services/data-api

RUN pnpm install --frozen-lockfile --filter data-api
RUN pnpm --filter data-api run build
RUN pnpm deploy --filter data-api --prod --legacy /app/deploy

# Stage 3 - Runtime
FROM node:22-alpine
WORKDIR /app
ENV NODE_ENV=production

RUN addgroup -S nodejs && adduser -S nestjs -G nodejs

COPY --from=builder /app/deploy ./
COPY --from=builder /app/services/data-api/dist ./dist

USER nestjs
EXPOSE 3000
CMD [ "node", "dist/main.js" ]
