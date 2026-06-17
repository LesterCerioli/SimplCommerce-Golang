FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/payment ./services/payment
RUN go build -o /app/payment ./services/payment/cmd/paymentd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/payment .
EXPOSE 8085
CMD ["./payment"]
