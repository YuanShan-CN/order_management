FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order_management .
RUN go build -o init_user scripts/init_user.go
RUN go build -o import_csv scripts/import_csv.go
RUN go build -o fill_shops scripts/fill_shops.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/order_management .
COPY --from=builder /app/init_user .
COPY --from=builder /app/import_csv .
COPY --from=builder /app/fill_shops .
COPY --from=builder /app/templates ./templates

EXPOSE 8080

CMD ["./order_management"]