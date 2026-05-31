# syntax=docker/dockerfile:1.7

# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org

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

# Download dependencies once and reuse them across parallel service builds.
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    sh -ec 'go work sync && for mod in ./pkg ./be/auth ./be/catalog ./be/chat ./be/gateway ./be/kitchen ./be/logistics ./be/notification ./be/orders ./be/treasury; do (cd "$mod" && go mod download); done'

# Copy the rest of the source code
COPY . .

# Build the application
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=0 GOOS=linux go build -o /app/bin/service ./be/treasury/cmd/main.go

# Stage 2: Final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/bin/service .
EXPOSE 9006
CMD ["./service"]
