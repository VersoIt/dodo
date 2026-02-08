package usecase

import (
	"context"
	"testing"

	"github.com/versoit/diploma/services/catalog"
)

type MockProductRepo struct {
	store map[string]*catalog.Product
}

func NewMockProductRepo() *MockProductRepo {
	return &MockProductRepo{store: make(map[string]*catalog.Product)}
}

func (m *MockProductRepo) Save(ctx context.Context, p *catalog.Product) error {
	m.store[p.ID()] = p
	return nil
}

func (m *MockProductRepo) FindByID(ctx context.Context, id string) (*catalog.Product, error) {
	if p, ok := m.store[id]; ok {
		return p, nil
	}
	return nil, catalog.ErrProductNotFound
}

func (m *MockProductRepo) FindAll(ctx context.Context) ([]*catalog.Product, error) {
	var list []*catalog.Product
	for _, p := range m.store {
		list = append(list, p)
	}
	return list, nil
}

func TestCatalogUseCase_CreateProduct(t *testing.T) {
	repo := NewMockProductRepo()
	uc := NewCatalogUseCase(repo)

	p, err := uc.CreateProduct(context.Background(), "Burger", "Delicious", catalog.CatClassic, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	saved, _ := repo.FindByID(context.Background(), p.ID())
	if saved == nil {
		t.Error("product not saved")
	}
}

func TestCatalogUseCase_UpdatePrice(t *testing.T) {
	repo := NewMockProductRepo()
	uc := NewCatalogUseCase(repo)
	p, _ := uc.CreateProduct(context.Background(), "Burger", "Desc", catalog.CatClassic, 100)

	err := uc.UpdatePrice(context.Background(), p.ID(), 150)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := repo.FindByID(context.Background(), p.ID())
	if updated.BasePrice() != 150 {
		t.Errorf("expected price 150, got %v", updated.BasePrice())
	}
}