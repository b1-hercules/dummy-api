# Stage 1: Build
FROM golang:1.23 AS builder
WORKDIR /app

# Copy semua file ke dalam container
COPY . .

# Download dependencies dan build aplikasi
RUN go mod tidy && go build -o main

# Stage 2: Deploy
FROM alpine:latest
WORKDIR /root/

# Copy binary dari stage build
COPY --from=builder /app/main .

# Jalankan aplikasi
CMD ["./main"]
