package usecase

import (
	"context"
	"fmt"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/treasury"
)

type TreasuryUseCase struct {
	repo treasury.PaymentRepository
}

func NewTreasuryUseCase(repo treasury.PaymentRepository) *TreasuryUseCase {
	return &TreasuryUseCase{repo: repo}
}

// InitiatePayment создает запись о намерении произвести платеж.
func (uc *TreasuryUseCase) InitiatePayment(ctx context.Context, orderID string, amount common.Money, method treasury.PaymentMethod) (*treasury.Payment, error) {
	payment := treasury.NewPayment(orderID, amount, method)
	
	if err := uc.repo.Save(payment); err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}
	
	return payment, nil
}

// ConfirmPayment подтверждает успешное прохождение транзакции от банка.
func (uc *TreasuryUseCase) ConfirmPayment(ctx context.Context, orderID string, transactionID string) error {
	payment, err := uc.repo.FindByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("payment record not found: %w", err)
	}

	if err := payment.Confirm(transactionID); err != nil {
		return err
	}

	return uc.repo.Save(payment)
}
