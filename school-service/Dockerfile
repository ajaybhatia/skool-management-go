FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy shared package
COPY shared ./shared

# Copy service code
COPY school-service ./school-service

# Build the application
WORKDIR /app/school-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/school-service/main .

EXPOSE 8082 50051

CMD ["./main"]
