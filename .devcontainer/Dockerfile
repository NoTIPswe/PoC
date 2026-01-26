FROM mcr.microsoft.com/devcontainers/base:ubuntu-22.04

ARG GO_VERSION=1.25.6
ARG TARGETARCH
ARG WORKSPACE_DIR=/workspaces/PoC

# Install curl and certificates
RUN apt-get update && apt-get install -y \
    curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install Go
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz \
    | tar -C /usr/local -xzf -

ENV PATH="/usr/local/go/bin:${PATH}"

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

RUN corepack enable

# Configure pnpm global bin path 
ENV PNPM_HOME="/home/vscode/.local/share/pnpm"
ENV PATH="${PNPM_HOME}:${PATH}"

# Install angular 
RUN pnpm install -g @angular/cli

# Set the working directory
RUN mkdir -p ${WORKSPACE_DIR} && chown vscode:vscode ${WORKSPACE_DIR}
WORKDIR ${WORKSPACE_DIR}

# Switch vscode user
USER vscode