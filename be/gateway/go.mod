module github.com/versoit/diploma/gateway

go 1.24.0

require (
	github.com/gofiber/contrib/websocket v1.3.4
	github.com/gofiber/fiber/v2 v2.52.11
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/versoit/diploma/be/auth v0.0.0-20260208211213-389e778a5129
	github.com/versoit/diploma/be/catalog v0.0.0-20260209221312-cf1acef61820
	github.com/versoit/diploma/be/chat v0.0.0-20260217214930-d3ab1be613b7
	github.com/versoit/diploma/be/kitchen v0.0.0-20260208211213-389e778a5129
	github.com/versoit/diploma/be/logistics v0.0.0-20260208211213-389e778a5129
	github.com/versoit/diploma/be/orders v0.0.0-20260208211213-389e778a5129
	github.com/versoit/diploma/pkg v0.0.0-20260214122137-9cedc0313251
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.65.0
	go.opentelemetry.io/otel/trace v1.40.0
	go.uber.org/fx v1.24.0
	google.golang.org/grpc v1.78.0
)

replace github.com/versoit/diploma/pkg => ../pkg

replace github.com/versoit/diploma/be/auth => ../auth

replace github.com/versoit/diploma/be/catalog => ../catalog

replace github.com/versoit/diploma/be/kitchen => ../kitchen

replace github.com/versoit/diploma/be/logistics => ../logistics

replace github.com/versoit/diploma/be/orders => ../orders

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2 v2.0.2 // indirect
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.2 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/fasthttp/websocket v1.5.8 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.7 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.8.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/gotils v0.0.0-20240303185622-093b76447511 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.52.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.40.0 // indirect
	go.opentelemetry.io/otel/metric v1.40.0 // indirect
	go.opentelemetry.io/otel/sdk v1.40.0 // indirect
	go.opentelemetry.io/proto/otlp v1.9.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260128011058-8636f8732409 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260128011058-8636f8732409 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/versoit/diploma/be/chat => ../chat
