package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/be/notification/internal/handler"
	"github.com/versoit/diploma/be/notification/internal/repository"
	"github.com/versoit/diploma/be/notification/internal/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		common.NewTransactionManager,
		repository.NewNotificationRepository,
		usecase.NewNotificationUseCase,
		handler.NewNotificationHandler,
	),
)
