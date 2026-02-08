package usecase

import (
	"github.com/versoit/diploma/services/treasury/internal/api/grpc"
	"github.com/versoit/diploma/services/treasury/internal/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryPaymentRepository,
		NewTreasuryUseCase,
		grpc.NewTreasuryHandler,
	),
)
