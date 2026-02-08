package usecase

import (
	"context"
	"fmt"

	"github.com/versoit/diploma/services/kitchen"
)

type KitchenUseCase struct {
	repo kitchen.TicketRepository
}

func NewKitchenUseCase(repo kitchen.TicketRepository) *KitchenUseCase {
	return &KitchenUseCase{repo: repo}
}

// AcceptOrder преобразует заказ в кухонный тикет.
func (uc *KitchenUseCase) AcceptOrder(ctx context.Context, orderID string, items []kitchen.KitchenItem) (*kitchen.KitchenTicket, error) {
	ticket := kitchen.NewTicket(orderID, items)

	if err := uc.repo.Save(ticket); err != nil {
		return nil, fmt.Errorf("failed to save kitchen ticket: %w", err)
	}

	return ticket, nil
}

// StartCooking переводит тикет в процесс приготовления.
func (uc *KitchenUseCase) StartCooking(ctx context.Context, ticketID string) error {
	ticket, err := uc.repo.FindByID(ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	if err := ticket.StartCooking(); err != nil {
		return err
	}

	return uc.repo.Save(ticket)
}

// MarkReady завершает приготовление заказа.
func (uc *KitchenUseCase) MarkReady(ctx context.Context, ticketID string) error {
	ticket, err := uc.repo.FindByID(ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	if err := ticket.MarkReady(); err != nil {
		return err
	}

	return uc.repo.Save(ticket)
}
