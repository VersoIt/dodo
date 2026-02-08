package usecase

import (
	"github.com/versoit/diploma/services/kitchen/internal/api/grpc"
	"github.com/versoit/diploma/services/kitchen/internal/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryTicketRepository,
		NewKitchenUseCase,
		grpc.NewKitchenHandler,
	),
)
