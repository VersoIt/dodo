package usecase

import (
	"github.com/versoit/diploma/services/orders/internal/api/grpc"
	"github.com/versoit/diploma/services/orders/internal/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryOrderRepository,
		NewOrderUseCase,
		grpc.NewOrdersHandler,
	),
)
