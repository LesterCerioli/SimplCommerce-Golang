FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/activitylog ./services/activitylog
RUN go build -o /app/activitylog ./services/activitylog/cmd/activitylogd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/activitylog .
EXPOSE 8092
CMD ["./activitylog"]
