# -------------------------------------------------------------
# Stage 1: Build the static Go binary
# -------------------------------------------------------------
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Install git and certificates in case private dependencies are needed
RUN apk add --no-cache git ca-certificates tzdata

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Compile fully static binary:
# - CGO_ENABLED=0 creates a pure Go binary with no dynamic libc dependencies
# - ldflags "-s -w" strips debugging symbols and DWARF tables to cut binary size in half
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/papercut \
    ./cmd/server/main.go

# -------------------------------------------------------------
# Stage 2: Minimal runtime image
# -------------------------------------------------------------
FROM alpine:3.20

# Add security updates, SSL certs, and standard timezones
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

# Copy only the compiled binary from builder
COPY --from=builder /bin/papercut /usr/local/bin/papercut

# Run as an unprivileged non-root user
USER appuser

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/papercut"]
