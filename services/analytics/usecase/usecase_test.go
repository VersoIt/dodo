package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/versoit/diploma/services/analytics"
)

type MockAnalyticsRepo struct {
	store map[string]*analytics.ManagerKPI
}

func (m *MockAnalyticsRepo) SaveKPI(ctx context.Context, k *analytics.ManagerKPI) error {
	m.store[k.ManagerID()] = k
	return nil
}

func (m *MockAnalyticsRepo) GetKPI(ctx context.Context, managerID string) (*analytics.ManagerKPI, error) {
	if k, ok := m.store[managerID]; ok {
		return k, nil
	}
	return nil, fmt.Errorf("not found")
}

func TestAnalyticsUseCase_RecordSale(t *testing.T) {
	repo := &MockAnalyticsRepo{store: make(map[string]*analytics.ManagerKPI)}
	uc := NewAnalyticsUseCase(repo)

	err := uc.RecordSale(context.Background(), "man1", 5000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	kpi, _ := repo.GetKPI(context.Background(), "man1")
	if kpi.Fact() != 5000 {
		t.Errorf("expected 5000, got %v", kpi.Fact())
	}
}
