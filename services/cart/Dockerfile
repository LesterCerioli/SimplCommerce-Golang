FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/cart ./services/cart
RUN go build -o /app/cart ./services/cart/cmd/cartd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/cart .
EXPOSE 8083
CMD ["./cart"]
