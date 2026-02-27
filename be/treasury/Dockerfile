# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

# First, copy over the files that define dependencies
COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY be/auth/go.mod be/auth/go.sum ./be/auth/
COPY be/catalog/go.mod be/catalog/go.sum ./be/catalog/
COPY be/chat/go.mod be/chat/go.sum ./be/chat/
COPY be/kitchen/go.mod be/kitchen/go.sum ./be/kitchen/
COPY be/logistics/go.mod be/logistics/go.sum ./be/logistics/
COPY be/notification/go.mod be/notification/go.sum ./be/notification/
COPY be/orders/go.mod be/orders/go.sum ./be/orders/
COPY be/treasury/go.mod be/treasury/go.sum ./be/treasury/
COPY be/gateway/go.mod be/gateway/go.sum ./be/gateway/

# Download dependencies. This layer is cached and only runs if go.mod/go.sum files change.
RUN go work sync

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/service ./be/treasury/cmd/main.go

# Stage 2: Final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/bin/service .
EXPOSE 9006
CMD ["./service"]
