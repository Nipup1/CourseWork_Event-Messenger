FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY config/container.yaml /root/config/local.yaml
EXPOSE 8080
CMD ["./main", "--config=/root/config/local.yaml"]