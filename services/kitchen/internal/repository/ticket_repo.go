package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/versoit/diploma/services/kitchen"
)

type InMemoryTicketRepository struct {
	mu    sync.RWMutex
	store map[string]*kitchen.KitchenTicket
}

func NewInMemoryTicketRepository() kitchen.TicketRepository {
	return &InMemoryTicketRepository{
		store: make(map[string]*kitchen.KitchenTicket),
	}
}

func (r *InMemoryTicketRepository) Save(ctx context.Context, t *kitchen.KitchenTicket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[t.ID()] = t
	return nil
}

func (r *InMemoryTicketRepository) FindByID(ctx context.Context, id string) (*kitchen.KitchenTicket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("ticket not found")
	}
	return t, nil
}

func (r *InMemoryTicketRepository) FindPending(ctx context.Context) ([]*kitchen.KitchenTicket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*kitchen.KitchenTicket
	for _, t := range r.store {
		if t.Status() == kitchen.TicketQueued {
			list = append(list, t)
		}
	}
	return list, nil
}
