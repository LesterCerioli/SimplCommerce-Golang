FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/api-gateway ./services/api-gateway
RUN go build -o /app/api-gateway ./services/api-gateway/cmd/gatewayd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]
