package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/versoit/diploma/be/chat/internal/domain"
)

type MessageService interface {
	Save(ctx context.Context, msg *domain.Message) error
	Broadcast(wsMsg domain.WSMessage)
	GetHistory(ctx context.Context, orderID uuid.UUID, limit int) ([]domain.Message, error)
	GetUpdates(ctx context.Context, orderID uuid.UUID, afterID int64) ([]domain.Message, error)
	MarkAsRead(ctx context.Context, msgID int64) error
}

type messageService struct {
	repo domain.MessageRepository
	hub  *Hub
	log  *slog.Logger
}

func NewMessageService(repo domain.MessageRepository, hub *Hub, log *slog.Logger) MessageService {
	return &messageService{
		repo: repo,
		hub:  hub,
		log:  log,
	}
}

func (s *messageService) Save(ctx context.Context, msg *domain.Message) error {
	if err := s.repo.Save(ctx, msg); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}

func (s *messageService) Broadcast(wsMsg domain.WSMessage) {
	s.hub.Broadcast <- wsMsg
}

func (s *messageService) GetHistory(ctx context.Context, orderID uuid.UUID, limit int) ([]domain.Message, error) {
	return s.repo.GetHistory(ctx, orderID, limit)
}

func (s *messageService) GetUpdates(ctx context.Context, orderID uuid.UUID, afterID int64) ([]domain.Message, error) {
	return s.repo.GetAfterID(ctx, orderID, afterID)
}

func (s *messageService) MarkAsRead(ctx context.Context, msgID int64) error {
	return s.repo.MarkAsRead(ctx, msgID)
}
