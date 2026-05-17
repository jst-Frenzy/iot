FROM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o wetSensor ./cmd/app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/wetSensor .
COPY config.yaml .

EXPOSE 8092

CMD ["./wetSensor"]
