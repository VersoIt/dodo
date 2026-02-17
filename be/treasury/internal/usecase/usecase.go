package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/be/treasury/internal/domain"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type TreasuryUseCase struct {
	repo domain.PaymentRepository
	tm   trm.Manager
	log  *slog.Logger
}

func NewTreasuryUseCase(repo domain.PaymentRepository, tm trm.Manager, log *slog.Logger) *TreasuryUseCase {
	return &TreasuryUseCase{
		repo: repo,
		tm:   tm,
		log:  log,
	}
}

func (uc *TreasuryUseCase) InitiatePayment(ctx context.Context, orderID string, amount common.Money, method domain.PaymentMethod) (*domain.Payment, error) {
	if orderID == "" {
		return nil, fmt.Errorf("%w: order ID is required", ErrInvalidInput)
	}
	if !amount.IsPositive() {
		return nil, fmt.Errorf("%w: payment amount must be positive", ErrInvalidInput)
	}

	uc.log.Info("initiating payment", slog.String("order_id", orderID), slog.Float64("amount", amount.InexactFloat64()))

	var payment *domain.Payment
	err := uc.tm.Do(ctx, func(ctx context.Context) error {
		payment = domain.NewPayment(orderID, amount, method)

		if err := uc.repo.Save(ctx, payment); err != nil {
			return fmt.Errorf("failed to register payment attempt for order %s: %w", orderID, err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *TreasuryUseCase) ConfirmPayment(ctx context.Context, orderID string, transactionID string) error {
	if orderID == "" || transactionID == "" {
		return fmt.Errorf("%w: order ID and transaction ID are required", ErrInvalidInput)
	}

	uc.log.Info("confirming payment", slog.String("order_id", orderID), slog.String("transaction_id", transactionID))

	return uc.tm.Do(ctx, func(ctx context.Context) error {
		payment, err := uc.repo.FindByOrderID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("find payment for order %s: %w", orderID, err)
		}

		if err := payment.Confirm(transactionID); err != nil {
			return fmt.Errorf("confirm payment: %w", err)
		}

		if err := uc.repo.Save(ctx, payment); err != nil {
			return fmt.Errorf("persist payment: %w", err)
		}
		return nil
	})
}
