package app

import (
	"github.com/versoit/diploma/services/treasury/internal/api/grpc"
	"github.com/versoit/diploma/services/treasury/internal/repository"
	"github.com/versoit/diploma/services/treasury/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryPaymentRepository,
		usecase.NewTreasuryUseCase,
		grpc.NewTreasuryHandler,
	),
)
