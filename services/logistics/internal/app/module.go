package app

import (
	"github.com/versoit/diploma/services/logistics/internal/api/grpc"
	"github.com/versoit/diploma/services/logistics/internal/repository"
	"github.com/versoit/diploma/services/logistics/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryDeliveryRepository,
		repository.NewInMemoryCourierRepository,
		usecase.NewLogisticsUseCase,
		grpc.NewLogisticsHandler,
	),
)
