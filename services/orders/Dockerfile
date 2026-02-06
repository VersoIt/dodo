# Stage 1: Build
FROM golang:1.25.3-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service from the cmd directory
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/service ./services/orders/cmd/main.go

# Stage 2: Final image
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/bin/service .

# Use environment variable for port flexibility
ENV PORT=8080
EXPOSE ${PORT}

CMD ["./service"]