package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/versoit/diploma/services/analytics"
)

type InMemoryAnalyticsRepository struct {
	mu    sync.RWMutex
	store map[string]*analytics.ManagerKPI
}

func NewInMemoryAnalyticsRepository() analytics.AnalyticsRepository {
	return &InMemoryAnalyticsRepository{
		store: make(map[string]*analytics.ManagerKPI),
	}
}

func (r *InMemoryAnalyticsRepository) SaveKPI(ctx context.Context, k *analytics.ManagerKPI) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[k.ManagerID()] = k
	return nil
}

func (r *InMemoryAnalyticsRepository) GetKPI(ctx context.Context, managerID string) (*analytics.ManagerKPI, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	k, ok := r.store[managerID]
	if !ok {
		return nil, fmt.Errorf("kpi not found")
	}
	return k, nil
}
