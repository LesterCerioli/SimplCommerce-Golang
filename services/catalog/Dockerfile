FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/catalog ./services/catalog
RUN go build -o /app/catalog ./services/catalog/cmd/catalogd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/catalog .
EXPOSE 8082
CMD ["./catalog"]
