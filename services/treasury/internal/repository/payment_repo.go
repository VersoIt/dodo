package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/versoit/diploma/services/treasury"
)

type InMemoryPaymentRepository struct {
	mu    sync.RWMutex
	store map[string]*treasury.Payment
}

func NewInMemoryPaymentRepository() treasury.PaymentRepository {
	return &InMemoryPaymentRepository{
		store: make(map[string]*treasury.Payment),
	}
}

func (r *InMemoryPaymentRepository) Save(ctx context.Context, p *treasury.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[p.OrderID()] = p
	return nil
}

func (r *InMemoryPaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*treasury.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[orderID]
	if !ok {
		return nil, fmt.Errorf("payment not found")
	}
	return p, nil
}
