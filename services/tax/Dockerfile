FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/tax ./services/tax
RUN go build -o /app/tax ./services/tax/cmd/taxd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/tax .
EXPOSE 8091
CMD ["./tax"]
