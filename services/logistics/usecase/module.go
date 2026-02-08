package usecase

import (
	"github.com/versoit/diploma/services/logistics"
	"github.com/versoit/diploma/services/logistics/internal/api/grpc"
	"github.com/versoit/diploma/services/logistics/internal/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryLogisticsRepository,
		// Fx can handle multiple return values, but let's be explicit if needed
		func(dr logistics.DeliveryRepository) logistics.DeliveryRepository { return dr },
		func(r *repository.InMemoryLogisticsRepository) logistics.CourierRepository { return r.ToCourierRepo() },
		NewLogisticsUseCase,
		grpc.NewLogisticsHandler,
	),
)
