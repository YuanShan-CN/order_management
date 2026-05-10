FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order_management .
RUN go build -o init_user scripts/init_user.go

FROM alpine:latest

WORKDIR /app

# 全部使用UTC时间，不需要时区配置
RUN mkdir -p /app/data/csv

COPY --from=builder /app/order_management .
COPY --from=builder /app/init_user .
COPY --from=builder /app/templates ./templates

ENV IS_DOCKER=true

EXPOSE 8080

CMD ["./order_management"]