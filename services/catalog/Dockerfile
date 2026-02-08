# Stage 1: Build
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/service ./services/catalog/cmd/main.go

# Stage 2: Final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/bin/service .
ENV PORT=8080
EXPOSE ${PORT}
CMD ["./service"]
