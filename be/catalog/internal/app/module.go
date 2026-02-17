package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/be/catalog/internal/handler"
	"github.com/versoit/diploma/be/catalog/internal/repository"
	usecase2 "github.com/versoit/diploma/be/catalog/internal/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		common.NewTransactionManager,
		repository.NewProductRepository,
		usecase2.NewCatalogUseCase,
		handler.NewCatalogHandler,
	),
)
