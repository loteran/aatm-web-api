# Build stage
FROM golang:1.22-bookworm AS builder

WORKDIR /build

# Copy go mod files first for caching
COPY api/go.mod ./
RUN go mod download 2>/dev/null || true

# Copy source code
COPY api/*.go ./
COPY api/static ./static/

# Download dependencies
RUN go mod tidy

# Build for ARM64 (Raspberry Pi)
RUN CGO_ENABLED=1 GOOS=linux go build -o aatm-api .

# Runtime stage
FROM debian:bookworm-slim

# Install dependencies (split for better QEMU compatibility)
RUN apt-get update && \
    apt-get install -y --no-install-recommends mediainfo supervisor && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get update && \
    apt-get install -y --no-install-recommends qbittorrent-nox && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates || true && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary
COPY --from=builder /build/aatm-api .

# Create directories
RUN mkdir -p /data /config/qBittorrent /downloads

# Copy supervisor config
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Copy qBittorrent default config
COPY qBittorrent.conf /config/qBittorrent/qBittorrent.conf

# Expose ports (API + qBittorrent WebUI)
EXPOSE 8080 8081

# Environment variables
ENV PORT=8080
ENV DATA_DIR=/data
ENV QBT_WEBUI_PORT=8081

# Run supervisor (manages both services)
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
