FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order_management .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/order_management .
COPY --from=builder /app/templates ./templates

EXPOSE 8080

CMD ["./order_management"]