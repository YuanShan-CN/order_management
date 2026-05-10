FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order_management .
RUN go build -o init_user scripts/init_user.go

FROM alpine:latest

WORKDIR /app

# 安装时区数据（用于定时任务计算配置的时区时间）
RUN apk add --no-cache tzdata
RUN mkdir -p /app/data/csv

COPY --from=builder /app/order_management .
COPY --from=builder /app/init_user .
COPY --from=builder /app/templates ./templates

ENV IS_DOCKER=true

EXPOSE 8080

CMD ["./order_management"]