FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o analysis-service ./services/analysis/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/analysis-service .
COPY --from=builder /app/services/subscription/migrations ./migrations
EXPOSE 5053

CMD ["./analysis-service"]