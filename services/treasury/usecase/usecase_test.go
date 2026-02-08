package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/versoit/diploma/services/treasury"
)

type MockTreasuryRepo struct {
	store map[string]*treasury.Payment
}

func (m *MockTreasuryRepo) Save(ctx context.Context, p *treasury.Payment) error {
	m.store[p.OrderID()] = p
	return nil
}
func (m *MockTreasuryRepo) FindByOrderID(ctx context.Context, id string) (*treasury.Payment, error) {
	if p, ok := m.store[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("not found")
}

func TestTreasuryUseCase_PaymentFlow(t *testing.T) {
	repo := &MockTreasuryRepo{store: make(map[string]*treasury.Payment)}
	uc := NewTreasuryUseCase(repo)

	_, err := uc.InitiatePayment(context.Background(), "ord-1", 1000, treasury.MethodCard)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	err = uc.ConfirmPayment(context.Background(), "ord-1", "trans-xyz")
	if err != nil {
		t.Fatalf("confirm failed: %v", err)
	}

	saved, _ := repo.FindByOrderID(context.Background(), "ord-1")
	if saved.Status() != treasury.PayStatusSuccess {
		t.Error("should be success")
	}
}
