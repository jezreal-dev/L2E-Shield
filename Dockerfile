# Build stage
FROM golang:1.26.4-alpine AS builder

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o l2e-shield cmd/server/main.go

# Run stage
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/l2e-shield .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./l2e-shield"]
