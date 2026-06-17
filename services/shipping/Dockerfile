FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/shipping ./services/shipping
RUN go build -o /app/shipping ./services/shipping/cmd/shippingd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/shipping .
EXPOSE 8086
CMD ["./shipping"]
