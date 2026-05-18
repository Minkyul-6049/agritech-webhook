# Step 1: Build stage (Matching local Go 1.22 version)
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Leverage Docker cache for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook-server main.go

# Step 2: Lightweight runtime stage
FROM alpine:3.19
WORKDIR /root/
RUN apk --no-cache add ca-certificates tzdata

# Copy built binary from builder stage
COPY --from=builder /app/webhook-server .

# 🚨 SECURITY REMOVAL: Do NOT copy .env inside the image.
# Environment variables will be injected dynamically at runtime via Docker/K3s.

EXPOSE 8080
CMD ["./webhook-server"]
