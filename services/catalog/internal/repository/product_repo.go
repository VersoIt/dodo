package repository

import (
	"context"
	"sync"

	"github.com/versoit/diploma/services/catalog"
)

type InMemoryProductRepository struct {
	mu    sync.RWMutex
	store map[string]*catalog.Product
}

func NewInMemoryProductRepository() catalog.ProductRepository {
	return &InMemoryProductRepository{
		store: make(map[string]*catalog.Product),
	}
}

func (r *InMemoryProductRepository) Save(ctx context.Context, p *catalog.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[p.ID()] = p
	return nil
}

func (r *InMemoryProductRepository) FindByID(ctx context.Context, id string) (*catalog.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[id]
	if !ok {
		return nil, catalog.ErrProductNotFound
	}
	return p, nil
}

func (r *InMemoryProductRepository) FindAll(ctx context.Context) ([]*catalog.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*catalog.Product, 0, len(r.store))
	for _, p := range r.store {
		list = append(list, p)
	}
	return list, nil
}
