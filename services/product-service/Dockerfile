FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy gen directory first
COPY gen ./gen

# Change workdir to match relative path structure
WORKDIR /app/services/product-service

# Copy go mod files
COPY services/product-service/go.mod services/product-service/go.sum ./
RUN go mod download

# Copy rest of the code
COPY services/product-service .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o product-service ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/services/product-service/product-service .
EXPOSE 50052

CMD ["./product-service"]