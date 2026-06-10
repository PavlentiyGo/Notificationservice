FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o currency-service ./services/currency/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/currency-service .

EXPOSE 5052

CMD ["./currency-service"]