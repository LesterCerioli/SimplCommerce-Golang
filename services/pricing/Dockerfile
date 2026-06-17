FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/pricing ./services/pricing
RUN go build -o /app/pricing ./services/pricing/cmd/pricingd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/pricing .
EXPOSE 8088
CMD ["./pricing"]
