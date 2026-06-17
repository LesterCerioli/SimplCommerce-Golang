FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/cms ./services/cms
RUN go build -o /app/cms ./services/cms/cmd/cmsd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/cms .
EXPOSE 8090
CMD ["./cms"]
