# Stage 1: Build
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

# Copy work files
COPY go.work go.work.sum* ./
COPY go.mod go.sum* ./

# Copy ALL module files for workspace sync
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY be/auth/go.mod be/auth/go.sum* ./be/auth/
COPY be/catalog/go.mod be/catalog/go.sum* ./be/catalog/
COPY be/kitchen/go.mod be/kitchen/go.sum* ./be/kitchen/
COPY be/logistics/go.mod be/logistics/go.sum* ./be/logistics/
COPY be/notification/go.mod be/notification/go.sum* ./be/notification/
COPY be/orders/go.mod be/orders/go.sum* ./be/orders/
COPY be/treasury/go.mod be/treasury/go.sum* ./be/treasury/
COPY be/gateway/go.mod be/gateway/go.sum* ./be/gateway/

ENV GOTOOLCHAIN=local
ENV GOTOOLCHAIN=local
RUN go work sync
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/service ./be/kitchen/cmd/main.go

# Stage 2: Final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/bin/service .
EXPOSE 8080
CMD ["./service"]



