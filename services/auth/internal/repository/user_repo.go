package repository

import (
	"context"
	"sync"

	"github.com/versoit/diploma/services/auth"
)

type InMemoryUserRepository struct {
	mu      sync.RWMutex
	byEmail map[string]*auth.User
	byID    map[string]*auth.User
}

func NewInMemoryUserRepository() auth.UserRepository {
	return &InMemoryUserRepository{
		byEmail: make(map[string]*auth.User),
		byID:    make(map[string]*auth.User),
	}
}

func (r *InMemoryUserRepository) Save(ctx context.Context, u *auth.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byEmail[u.Email()] = u
	r.byID[u.ID()] = u
	return nil
}

func (r *InMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byEmail[email]
	if !ok {
		return nil, auth.ErrUserNotFound
	}
	return u, nil
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id string) (*auth.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, auth.ErrUserNotFound
	}
	return u, nil
}
