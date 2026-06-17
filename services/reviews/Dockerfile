FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.work go.work.sum ./
COPY pkg ./pkg
COPY services/reviews ./services/reviews
RUN go build -o /app/reviews ./services/reviews/cmd/reviewsd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/reviews .
EXPOSE 8089
CMD ["./reviews"]
