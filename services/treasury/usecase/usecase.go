package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/treasury"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type TreasuryUseCase struct {
	repo treasury.PaymentRepository
	log  *slog.Logger
}

func NewTreasuryUseCase(repo treasury.PaymentRepository, log *slog.Logger) *TreasuryUseCase {
	return &TreasuryUseCase{
		repo: repo,
		log:  log,
	}
}

func (uc *TreasuryUseCase) InitiatePayment(ctx context.Context, orderID string, amount common.Money, method treasury.PaymentMethod) (*treasury.Payment, error) {
	if orderID == "" {
		return nil, fmt.Errorf("%w: order ID is required", ErrInvalidInput)
	}
	if !amount.IsPositive() {
		return nil, fmt.Errorf("%w: payment amount must be positive", ErrInvalidInput)
	}

	uc.log.Info("initiating payment", slog.String("order_id", orderID), slog.Float64("amount", amount.InexactFloat64()))

	payment := treasury.NewPayment(orderID, amount, method)

	if err := uc.repo.Save(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to register payment attempt for order %s: %w", orderID, err)
	}

	return payment, nil
}

func (uc *TreasuryUseCase) ConfirmPayment(ctx context.Context, orderID string, transactionID string) error {
	if orderID == "" || transactionID == "" {
		return fmt.Errorf("%w: order ID and transaction ID are required", ErrInvalidInput)
	}

	uc.log.Info("confirming payment", slog.String("order_id", orderID), slog.String("transaction_id", transactionID))

	payment, err := uc.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("payment record for order %s not found: %w", orderID, err)
	}

	if err := payment.Confirm(transactionID); err != nil {
		return fmt.Errorf("domain logic error while confirming payment: %w", err)
	}

	if err := uc.repo.Save(ctx, payment); err != nil {
		return fmt.Errorf("failed to persist payment confirmation: %w", err)
	}

	return nil
}
