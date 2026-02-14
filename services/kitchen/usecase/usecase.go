package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/versoit/diploma/services/kitchen"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type KitchenUseCase struct {
	repo kitchen.TicketRepository
	log  *slog.Logger
}

func NewKitchenUseCase(repo kitchen.TicketRepository, log *slog.Logger) *KitchenUseCase {
	return &KitchenUseCase{
		repo: repo,
		log:  log,
	}
}

func (uc *KitchenUseCase) AcceptOrder(ctx context.Context, orderID string, items []kitchen.KitchenItem) (*kitchen.KitchenTicket, error) {
	if orderID == "" {
		return nil, fmt.Errorf("%w: order ID is required", ErrInvalidInput)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("%w: ticket must contain items", ErrInvalidInput)
	}

	ticket := kitchen.NewTicket(orderID, items)

	uc.log.Info("accepting order in kitchen", slog.String("order_id", orderID), slog.String("ticket_id", ticket.ID()))

	if err := uc.repo.Save(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to create kitchen ticket: %w", err)
	}

	return ticket, nil
}

func (uc *KitchenUseCase) StartCooking(ctx context.Context, ticketID string) error {
	if ticketID == "" {
		return fmt.Errorf("%w: ticket ID is required", ErrInvalidInput)
	}

	ticket, err := uc.repo.FindByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to find kitchen ticket %s: %w", ticketID, err)
	}

	if err := ticket.StartCooking(); err != nil {
		return fmt.Errorf("could not start cooking for ticket %s: %w", ticketID, err)
	}

	if err := uc.repo.Save(ctx, ticket); err != nil {
		return fmt.Errorf("failed to update ticket status to cooking: %w", err)
	}

	uc.log.Info("started cooking", slog.String("ticket_id", ticketID), slog.String("order_id", ticket.OrderID()))
	return nil
}

func (uc *KitchenUseCase) MarkReady(ctx context.Context, ticketID string) error {
	if ticketID == "" {
		return fmt.Errorf("%w: ticket ID is required", ErrInvalidInput)
	}

	ticket, err := uc.repo.FindByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to find kitchen ticket %s: %w", ticketID, err)
	}

	if err := ticket.MarkReady(); err != nil {
		return fmt.Errorf("could not mark ticket %s as ready: %w", ticketID, err)
	}

	if err := uc.repo.Save(ctx, ticket); err != nil {
		return fmt.Errorf("failed to update ticket status to ready: %w", err)
	}

	uc.log.Info("ticket ready", slog.String("ticket_id", ticketID), slog.String("order_id", ticket.OrderID()))
	return nil
}
