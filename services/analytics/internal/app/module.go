package app

import (
	"github.com/versoit/diploma/services/analytics/internal/api/grpc"
	"github.com/versoit/diploma/services/analytics/internal/repository"
	"github.com/versoit/diploma/services/analytics/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryAnalyticsRepository,
		usecase.NewAnalyticsUseCase,
		grpc.NewAnalyticsHandler,
	),
)
