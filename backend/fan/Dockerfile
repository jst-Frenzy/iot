FROM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o fan ./cmd/app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/fan .
COPY config.yaml .

EXPOSE 8093

CMD ["./fan"]
