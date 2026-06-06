FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o subscription-service ./services/subscription/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/subscription-service .
COPY --from=builder /app/services/subscription/migrations ./migrations

EXPOSE 5051

CMD ["./subscription-service"]