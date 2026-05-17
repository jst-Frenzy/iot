FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o temperatureSensor ./cmd/app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/temperatureSensor .
COPY config.yaml .

EXPOSE 8091

CMD ["./temperatureSensor"]
