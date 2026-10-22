# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy module files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with CGO disabled for alpine static compilation
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server/main.go

# Production stage
FROM alpine:3.20

WORKDIR /app

# Install ca-certificates and tzdata for TLS and timezone support
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create logs directory with proper ownership
RUN mkdir -p /app/logs && chown -R appuser:appgroup /app

# Copy compiled binary from builder stage
COPY --from=builder --chown=appuser:appgroup /app/server /app/server

USER appuser

EXPOSE 3000

CMD ["/app/server"]
