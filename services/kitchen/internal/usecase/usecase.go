package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/services/kitchen/internal/domain"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type KitchenUseCase struct {
	repo domain.TicketRepository
	tm   trm.Manager
	log  *slog.Logger
}

func NewKitchenUseCase(repo domain.TicketRepository, tm trm.Manager, log *slog.Logger) *KitchenUseCase {
	return &KitchenUseCase{
		repo: repo,
		tm:   tm,
		log:  log,
	}
}

func (uc *KitchenUseCase) AcceptOrder(ctx context.Context, orderID string, orderNumber string, items []domain.KitchenItem) (*domain.KitchenTicket, error) {
	if orderID == "" {
		return nil, fmt.Errorf("%w: order ID is required", ErrInvalidInput)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("%w: ticket must contain items", ErrInvalidInput)
	}

	var ticket *domain.KitchenTicket
	err := uc.tm.Do(ctx, func(ctx context.Context) error {
		ticket = domain.NewTicket(orderID, orderNumber, items)

		uc.log.Info("accepting order in kitchen", slog.String("order_id", orderID), slog.String("order_number", orderNumber), slog.String("ticket_id", ticket.ID()))

		if err := uc.repo.Save(ctx, ticket); err != nil {
			return fmt.Errorf("failed to create kitchen ticket: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (uc *KitchenUseCase) ListTickets(ctx context.Context) ([]*domain.KitchenTicket, error) {
	tickets, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list kitchen tickets: %w", err)
	}
	return tickets, nil
}

func (uc *KitchenUseCase) StartCooking(ctx context.Context, ticketID string) (string, error) {
	if ticketID == "" {
		return "", fmt.Errorf("%w: ticket ID is required", ErrInvalidInput)
	}

	var orderID string
	err := uc.tm.Do(ctx, func(ctx context.Context) error {
		ticket, err := uc.repo.FindByID(ctx, ticketID)
		if err != nil {
			return fmt.Errorf("find kitchen ticket %s: %w", ticketID, err)
		}

		if err := ticket.StartCooking(); err != nil {
			return fmt.Errorf("start cooking for ticket %s: %w", ticketID, err)
		}

		if err := uc.repo.Save(ctx, ticket); err != nil {
			return fmt.Errorf("update ticket status to cooking: %w", err)
		}

		orderID = ticket.OrderID()
		uc.log.Info("started cooking", slog.String("ticket_id", ticketID), slog.String("order_id", orderID))
		return nil
	})

	return orderID, err
}

func (uc *KitchenUseCase) MarkReady(ctx context.Context, ticketID string) (string, error) {
	if ticketID == "" {
		return "", fmt.Errorf("%w: ticket ID is required", ErrInvalidInput)
	}

	var orderID string
	err := uc.tm.Do(ctx, func(ctx context.Context) error {
		ticket, err := uc.repo.FindByID(ctx, ticketID)
		if err != nil {
			return fmt.Errorf("find kitchen ticket %s: %w", ticketID, err)
		}

		if err := ticket.MarkReady(); err != nil {
			return fmt.Errorf("mark ticket %s as ready: %w", ticketID, err)
		}

		if err := uc.repo.Save(ctx, ticket); err != nil {
			return fmt.Errorf("update ticket status to ready: %w", err)
		}

		orderID = ticket.OrderID()
		uc.log.Info("ticket ready", slog.String("ticket_id", ticketID), slog.String("order_id", orderID))
		return nil
	})

	return orderID, err
}
