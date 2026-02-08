package app

import (
	"github.com/versoit/diploma/services/orders/internal/api/grpc"
	"github.com/versoit/diploma/services/orders/internal/repository"
	"github.com/versoit/diploma/services/orders/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryOrderRepository,
		usecase.NewOrderUseCase,
		grpc.NewOrdersHandler,
	),
)
