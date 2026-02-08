package usecase

import (
	"github.com/versoit/diploma/services/catalog/internal/api/grpc"
	"github.com/versoit/diploma/services/catalog/internal/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryProductRepository,
		NewCatalogUseCase,
		grpc.NewCatalogHandler,
	),
)
