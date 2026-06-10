FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api-gateway ./services/api-gateway/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/api-gateway .

EXPOSE 5050

CMD ["./api-gateway"]