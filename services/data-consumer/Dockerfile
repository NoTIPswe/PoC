FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o data-consumer main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/data-consumer .

CMD ["./data-consumer"]