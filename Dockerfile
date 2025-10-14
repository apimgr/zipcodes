# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildDate=${BUILD_DATE} -w -s" \
    -o zipcodes \
    ./src

# Runtime stage
FROM alpine:latest

# Add ca-certificates
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -S zipcodes && adduser -S zipcodes -G zipcodes

# Copy binary
COPY --from=builder /build/zipcodes /usr/local/bin/zipcodes

# Create directories
RUN mkdir -p /config /data /logs && \
    chown -R zipcodes:zipcodes /config /data /logs

# Switch to non-root
USER zipcodes

# Set working directory
WORKDIR /data

# Metadata labels
LABEL org.opencontainers.image.source="https://github.com/apimgr/zipcodes"
LABEL org.opencontainers.image.description="zipcodes server"
LABEL org.opencontainers.image.licenses="MIT"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["zipcodes", "--status"]

EXPOSE 8080

ENTRYPOINT ["zipcodes"]
