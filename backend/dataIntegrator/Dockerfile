FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o dataIntegrator ./cmd/app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/dataIntegrator .
COPY config.yaml .
COPY credentials.yaml .
COPY migrations ./migrations

EXPOSE 8090 8099

CMD ["./dataIntegrator"]
