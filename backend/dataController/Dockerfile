FROM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o dataController ./cmd/app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/dataController .
COPY config.yaml .
COPY credentials.yaml .

EXPOSE 8095

CMD ["./dataController"]
