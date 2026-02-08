package app

import (
	"github.com/versoit/diploma/services/notification/internal/api/grpc"
	"github.com/versoit/diploma/services/notification/internal/repository"
	"github.com/versoit/diploma/services/notification/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryNotificationRepository,
		usecase.NewNotificationUseCase,
		grpc.NewNotificationHandler,
	),
)
