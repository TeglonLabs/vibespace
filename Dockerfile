# Multi-stage Dockerfile for production-ready Go application
# Stage 1: Build environment
FROM golang:1.24-alpine AS builder

# Build arguments for version information
ARG VERSION=dev
ARG BUILD_TIME
ARG COMMIT=unknown

# Set necessary environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -s /bin/sh -u 1001 appuser

# Set working directory
WORKDIR /build

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations and version info
RUN go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${COMMIT}" \
    -o vibespace-mcp \
    ./cmd/server

# Verify the binary
RUN ./vibespace-mcp --version || echo "Binary built successfully"

# Stage 2: Runtime environment
FROM scratch

# Import timezone data and certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary from builder
COPY --from=builder /build/vibespace-mcp /vibespace-mcp

# Use non-root user
USER appuser

# Expose default port (if applicable)
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/vibespace-mcp", "--health-check"] || exit 1

# Set default command
ENTRYPOINT ["/vibespace-mcp"]
CMD ["--help"]

# Metadata
LABEL maintainer="Team Topos <team@topos.com>" \
      org.opencontainers.image.title="vibespace-mcp" \
      org.opencontainers.image.description="MCP (Model Context Protocol) experience for managing vibes and worlds" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.source="https://github.com/bmorphism/vibespace-mcp-go" \
      org.opencontainers.image.vendor="Team Topos" \
      org.opencontainers.image.licenses="MIT"
