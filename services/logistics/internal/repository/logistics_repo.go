package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/versoit/diploma/services/logistics"
)

type InMemoryDeliveryRepository struct {
	mu    sync.RWMutex
	store map[string]*logistics.Delivery
}

func NewInMemoryDeliveryRepository() logistics.DeliveryRepository {
	return &InMemoryDeliveryRepository{
		store: make(map[string]*logistics.Delivery),
	}
}

func (r *InMemoryDeliveryRepository) Save(ctx context.Context, d *logistics.Delivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[d.OrderID()] = d
	return nil
}

func (r *InMemoryDeliveryRepository) FindByOrderID(ctx context.Context, id string) (*logistics.Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("delivery not found")
	}
	return d, nil
}

type InMemoryCourierRepository struct {
	mu    sync.RWMutex
	store map[string]*logistics.Courier
}

func NewInMemoryCourierRepository() logistics.CourierRepository {
	return &InMemoryCourierRepository{
		store: make(map[string]*logistics.Courier),
	}
}

func (r *InMemoryCourierRepository) Save(ctx context.Context, c *logistics.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[c.ID()] = c
	return nil
}

func (r *InMemoryCourierRepository) FindAvailable(ctx context.Context) ([]*logistics.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*logistics.Courier
	for _, c := range r.store {
		if c.Status() == logistics.CourierFree {
			list = append(list, c)
		}
	}
	return list, nil
}

func (r *InMemoryCourierRepository) FindByID(ctx context.Context, id string) (*logistics.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("courier not found")
	}
	return c, nil
}

func (r *InMemoryCourierRepository) UpdateLocation(ctx context.Context, id string, lat, lng float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.store[id]; ok {
		return c.UpdateLocation(lat, lng)
	}
	return fmt.Errorf("courier not found")
}