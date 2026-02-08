package app

import (
	"github.com/versoit/diploma/services/catalog/internal/api/grpc"
	"github.com/versoit/diploma/services/catalog/internal/repository"
	"github.com/versoit/diploma/services/catalog/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryProductRepository,
		usecase.NewCatalogUseCase,
		grpc.NewCatalogHandler,
	),
)
