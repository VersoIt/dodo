package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/catalog/internal/api/grpc"
	"github.com/versoit/diploma/services/catalog/internal/repository"
	"github.com/versoit/diploma/services/catalog/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		repository.NewProductRepository,
		usecase.NewCatalogUseCase,
		grpc.NewCatalogHandler,
	),
)
