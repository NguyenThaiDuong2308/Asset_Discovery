FROM golang:1.18-alpine AS builder

WORKDIR /app

# Cài đặt dependencies
RUN apk add --no-cache git

# Copy Go module files
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build ứng dụng
RUN CGO_ENABLED=0 GOOS=linux go build -o log-analysis-service ./cmd/main.go

# Tạo final image
FROM alpine:3.16

WORKDIR /app

# Cài đặt các tools cần thiết
RUN apk add --no-cache ca-certificates tzdata

# Copy binary từ builder stage
COPY --from=builder /app/log-analysis-service .
# Set timezone
ENV TZ=UTC

# Expose API port
EXPOSE 8080

# Run the application
CMD ["./log-analysis-service"]