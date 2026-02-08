package repository

import (
	"context"
	"sync"

	"github.com/versoit/diploma/services/orders"
)

type InMemoryOrderRepository struct {
	mu    sync.RWMutex
	store map[string]*orders.Order
}

func NewInMemoryOrderRepository() orders.OrderRepository {
	return &InMemoryOrderRepository{
		store: make(map[string]*orders.Order),
	}
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, o *orders.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[o.ID()] = o
	return nil
}

func (r *InMemoryOrderRepository) FindByID(ctx context.Context, id string) (*orders.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.store[id]
	if !ok {
		return nil, orders.ErrOrderNotFound
	}
	return o, nil
}
