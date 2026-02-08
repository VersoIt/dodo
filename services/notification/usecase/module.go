package usecase

import (
	"github.com/versoit/diploma/services/notification/internal/api/grpc"
	"github.com/versoit/diploma/services/notification/internal/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryNotificationRepository,
		NewNotificationUseCase,
		grpc.NewNotificationHandler,
	),
)
