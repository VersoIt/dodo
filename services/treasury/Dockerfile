# Stage 1: Build
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

# Copy work files
COPY go.work go.work.sum* ./
COPY go.mod go.sum* ./

# Copy ALL module files for workspace sync
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY services/auth/go.mod services/auth/go.sum* ./services/auth/
COPY services/catalog/go.mod services/catalog/go.sum* ./services/catalog/
COPY services/kitchen/go.mod services/kitchen/go.sum* ./services/kitchen/
COPY services/logistics/go.mod services/logistics/go.sum* ./services/logistics/
COPY services/notification/go.mod services/notification/go.sum* ./services/notification/
COPY services/orders/go.mod services/orders/go.sum* ./services/orders/
COPY services/treasury/go.mod services/treasury/go.sum* ./services/treasury/
COPY gateway/go.mod gateway/go.sum* ./gateway/

RUN go work sync
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/service ./services/treasury/cmd/main.go

# Stage 2: Final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/bin/service .
EXPOSE 8080
CMD ["./service"]



