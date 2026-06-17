FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/inventory ./services/inventory
RUN go build -o /app/inventory ./services/inventory/cmd/inventoryd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/inventory .
EXPOSE 8087
CMD ["./inventory"]
