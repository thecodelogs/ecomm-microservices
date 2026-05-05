FROM golang:1.25-alpine AS builder

WORKDIR /app

# ✅ COPY GEN FIRST (IMPORTANT)
COPY gen ./gen

# Change workdir to match the relative path structure expected by go.mod
WORKDIR /app/services/user-service

# Copy go mod files
COPY services/user-service/go.mod services/user-service/go.sum ./

# Now dependencies will work since ../../gen resolves to /app/gen
RUN go mod download

# Copy rest of code
COPY services/user-service .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service main.go
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Final stage
FROM alpine:latest
WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/services/user-service/user-service .
COPY --from=builder /go/bin/grpcurl /usr/local/bin/

CMD ["./user-service"]