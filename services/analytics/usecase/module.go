package usecase

import (
	"github.com/versoit/diploma/services/analytics/internal/api/grpc"
	"github.com/versoit/diploma/services/analytics/internal/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryAnalyticsRepository,
		NewAnalyticsUseCase,
		grpc.NewAnalyticsHandler,
	),
)
