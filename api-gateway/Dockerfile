FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy shared package
COPY shared ./shared

# Copy service code
COPY api-gateway ./api-gateway

# Build the application
WORKDIR /app/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/api-gateway/main .

EXPOSE 8080

CMD ["./main"]
