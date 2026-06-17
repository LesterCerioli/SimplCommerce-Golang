FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/orders ./services/orders
RUN go build -o /app/orders ./services/orders/cmd/ordersd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/orders .
EXPOSE 8084
CMD ["./orders"]
