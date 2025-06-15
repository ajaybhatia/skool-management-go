FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy shared package
COPY shared ./shared

# Copy service code
COPY auth-service ./auth-service

# Build the application
WORKDIR /app/auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/auth-service/main .

EXPOSE 8081

CMD ["./main"]
