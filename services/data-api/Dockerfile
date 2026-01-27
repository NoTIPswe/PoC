# Stage 1
FROM node:22-alpine AS base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable 
WORKDIR /app

# Stage 2 - dependencies
FROM base AS deps 
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Stage 3 - buildva
FROM base AS builder 
COPY --from=deps /app/node_modules ./node_modules
COPY . . 
RUN pnpm run build 
RUN pnpm prune --prod

# Stage 4 - final image
FROM node:22-alpine
WORKDIR /app 
ENV NODE_ENV=production 

RUN addgroup -S nodejs && adduser -S nestjs -G nodejs

# Copy only necessary things 
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

USER nestjs

EXPOSE 3000
CMD [ "node", "dist/main.js" ]
