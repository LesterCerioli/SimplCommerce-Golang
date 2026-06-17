FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/notifications ./services/notifications
RUN go build -o /app/notifications ./services/notifications/cmd/notificationsd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/notifications .
EXPOSE 8093
CMD ["./notifications"]
