package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/be/catalog/internal/domain"
)

type MockProductRepo struct {
	store map[string]*domain.Product
}

func NewMockProductRepo() *MockProductRepo {
	return &MockProductRepo{store: make(map[string]*domain.Product)}
}

func (m *MockProductRepo) Save(ctx context.Context, p *domain.Product) error {
	m.store[p.ID()] = p
	return nil
}

func (m *MockProductRepo) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	if p, ok := m.store[id]; ok {
		return p, nil
	}
	return nil, domain.ErrProductNotFound
}

func (m *MockProductRepo) FindAll(ctx context.Context) ([]*domain.Product, error) {
	var list []*domain.Product
	for _, p := range m.store {
		list = append(list, p)
	}
	return list, nil
}

type dummyTM struct{}

func (d *dummyTM) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (d *dummyTM) DoWithSettings(ctx context.Context, _ trm.Settings, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestCatalogUseCase_CreateProduct(t *testing.T) {
	repo := NewMockProductRepo()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewCatalogUseCase(repo, &dummyTM{}, log)

	p, err := uc.CreateProduct(context.Background(), "Burger", "Delicious", domain.CatClassic, 100, "http://img.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	saved, err := repo.FindByID(context.Background(), p.ID())
	if err != nil {
		t.Fatalf("failed to find product: %v", err)
	}
	if saved == nil {
		t.Error("product not saved")
	}
}

func TestCatalogUseCase_UpdatePrice(t *testing.T) {
	repo := NewMockProductRepo()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewCatalogUseCase(repo, &dummyTM{}, log)
	p, err := uc.CreateProduct(context.Background(), "Burger", "Desc", domain.CatClassic, 100, "")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	err = uc.UpdatePrice(context.Background(), p.ID(), 150)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := repo.FindByID(context.Background(), p.ID())
	if err != nil {
		t.Fatalf("failed to find updated product: %v", err)
	}
	if updated.BasePrice().InexactFloat64() != 150 {
		t.Errorf("expected price 150, got %v", updated.BasePrice())
	}
}
