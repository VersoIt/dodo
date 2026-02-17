package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/be/treasury/internal/domain"
)

type MockTreasuryRepo struct {
	store map[string]*domain.Payment
}

func (m *MockTreasuryRepo) Save(ctx context.Context, p *domain.Payment) error {
	m.store[p.OrderID()] = p
	return nil
}
func (m *MockTreasuryRepo) FindByOrderID(ctx context.Context, id string) (*domain.Payment, error) {
	if p, ok := m.store[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("not found")
}

type dummyTM struct{}

func (d *dummyTM) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (d *dummyTM) DoWithSettings(ctx context.Context, _ trm.Settings, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestTreasuryUseCase_PaymentFlow(t *testing.T) {
	repo := &MockTreasuryRepo{store: make(map[string]*domain.Payment)}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewTreasuryUseCase(repo, &dummyTM{}, log)

	_, err := uc.InitiatePayment(context.Background(), "ord-1", common.NewMoney(1000), domain.MethodCard)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	err = uc.ConfirmPayment(context.Background(), "ord-1", "trans-xyz")
	if err != nil {
		t.Fatalf("confirm failed: %v", err)
	}

	saved, err := repo.FindByOrderID(context.Background(), "ord-1")
	if err != nil {
		t.Fatalf("failed to find payment: %v", err)
	}
	if saved.Status() != domain.PayStatusSuccess {
		t.Error("should be success")
	}
}
