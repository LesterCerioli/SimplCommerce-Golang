FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/search ./services/search
RUN go build -o /app/search ./services/search/cmd/searchd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/search .
EXPOSE 8094
CMD ["./search"]
