package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/versoit/diploma/services/logistics"
)

type InMemoryLogisticsRepository struct {
	mu         sync.RWMutex
	deliveries map[string]*logistics.Delivery
	couriers   map[string]*logistics.Courier
}

func NewInMemoryLogisticsRepository() (logistics.DeliveryRepository, logistics.CourierRepository) {
	repo := &InMemoryLogisticsRepository{
		deliveries: make(map[string]*logistics.Delivery),
		couriers:   make(map[string]*logistics.Courier),
	}
	return repo, repo
}

// DeliveryRepo implementation
func (r *InMemoryLogisticsRepository) Save(ctx context.Context, d *logistics.Delivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deliveries[d.OrderID()] = d
	return nil
}

func (r *InMemoryLogisticsRepository) FindByOrderID(ctx context.Context, id string) (*logistics.Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.deliveries[id]
	if !ok {
		return nil, fmt.Errorf("delivery not found")
	}
	return d, nil
}

// CourierRepo implementation
func (r *InMemoryLogisticsRepository) SaveCourier(ctx context.Context, c *logistics.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.couriers[c.ID()] = c
	return nil
}

func (r *InMemoryLogisticsRepository) FindAvailable(ctx context.Context) ([]*logistics.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*logistics.Courier
	for _, c := range r.couriers {
		if c.Status() == logistics.CourierFree {
			list = append(list, c)
		}
	}
	return list, nil
}

func (r *InMemoryLogisticsRepository) FindByID(ctx context.Context, id string) (*logistics.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.couriers[id]
	if !ok {
		return nil, fmt.Errorf("courier not found")
	}
	return c, nil
}

func (r *InMemoryLogisticsRepository) UpdateLocation(ctx context.Context, id string, lat, lng float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.couriers[id]; ok {
		return c.UpdateLocation(lat, lng)
	}
	return fmt.Errorf("courier not found")
}

// Wrapper to match interface expectation if fx gets confused with dual returns
func (r *InMemoryLogisticsRepository) ToCourierRepo() logistics.CourierRepository {
	return &courierRepoWrapper{r}
}

type courierRepoWrapper struct{ *InMemoryLogisticsRepository }
func (w *courierRepoWrapper) Save(ctx context.Context, c *logistics.Courier) error { return w.SaveCourier(ctx, c) }
