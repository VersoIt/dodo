FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the entire project context for workspace support
COPY . .

# Build the gateway
RUN go build -o /app/gateway_bin ./gateway/cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/gateway_bin .

EXPOSE 8080

CMD ["./gateway_bin"]
