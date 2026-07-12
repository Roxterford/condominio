FROM ubuntu:24.04

# Versions
ARG GO_VERSION=1.26.0
ARG NODE_VERSION=24

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    wget \
    git \
    ca-certificates \
    gnupg \
    lsb-release \
    unzip \
    build-essential \
    sudo \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m devuser -s /bin/bash \
    && echo "devuser ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/devuser \
    && chmod 0440 /etc/sudoers.d/devuser

# Install Node.js (via NodeSource)
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# Install Go system-wide
RUN wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz

# Switch to devuser for user-level installs
USER devuser

# Install Bun
RUN curl -fsSL https://bun.sh/install | bash

# Install moon
RUN curl -fsSL https://moonrepo.dev/install/moon.sh | bash

# Set up PATH
ENV GOPATH="/home/devuser/go"
ENV PATH="/usr/local/go/bin:${GOPATH}/bin:/home/devuser/.bun/bin:/home/devuser/.moon/bin:${PATH}"

# Install Go tools
RUN go install github.com/segmentio/golines@latest
RUN go install github.com/99designs/gqlgen@latest
RUN go install github.com/evilmartians/lefthook/v2@latest

# Back to root for npm global install
USER root
RUN npm install -g @commitlint/cli commitizen

# Back to devuser for workspace
USER devuser
WORKDIR /home/devuser/workspace
