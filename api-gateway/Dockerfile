FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy gen directory first
COPY gen ./gen

# Change workdir to api-gateway
WORKDIR /app/api-gateway

# Copy go mod files from api-gateway directory
COPY api-gateway/go.mod api-gateway/go.sum ./
RUN go mod download

# Copy rest of the api-gateway code
COPY api-gateway .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/api-gateway/api-gateway .
COPY api-gateway/swagger.yaml ./swagger.yaml
EXPOSE 8080

CMD ["./api-gateway"]