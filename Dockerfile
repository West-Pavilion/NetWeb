# Multi-stage build for NetWeb

# Stage 1: Build Frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o netweb-server .

# Stage 3: Final image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl iputils bind-tools traceroute

WORKDIR /root/

# Copy backend binary
COPY --from=backend-builder /app/netweb-server .

# Copy frontend build
COPY --from=frontend-builder /app/frontend/build ./frontend/build

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./netweb-server"]