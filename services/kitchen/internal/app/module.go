package app

import (
	"github.com/versoit/diploma/services/kitchen/internal/api/grpc"
	"github.com/versoit/diploma/services/kitchen/internal/repository"
	"github.com/versoit/diploma/services/kitchen/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryTicketRepository,
		usecase.NewKitchenUseCase,
		grpc.NewKitchenHandler,
	),
)
