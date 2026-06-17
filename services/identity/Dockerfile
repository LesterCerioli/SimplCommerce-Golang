FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/identity ./services/identity
RUN go build -o /app/identity ./services/identity/cmd/identityd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/identity .
EXPOSE 8081
CMD ["./identity"]
