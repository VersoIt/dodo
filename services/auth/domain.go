package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// --- Value Objects & Enums ---

type Role int

const (
	RoleAdmin   Role = 0
	RoleManager Role = 1
	RoleChef    Role = 2
	RoleCourier Role = 3
	RoleClient  Role = 4
)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleManager:
		return "manager"
	case RoleChef:
		return "chef"
	case RoleCourier:
		return "courier"
	case RoleClient:
		return "client"
	default:
		return "unknown"
	}
}

// --- Errors ---

var (
	ErrInvalidEmail      = errors.New("invalid email")
	ErrWeakPassword      = errors.New("password must be at least 6 chars")
	ErrInsufficientBonus = errors.New("not enough bonus points")
	ErrUserNotFound      = errors.New("user not found")
)

// --- Aggregate ---

type User struct {
	id           string
	email        string
	passwordHash string
	role         Role
	createdAt    time.Time
	updatedAt    time.Time

	isClient    bool
	name        string
	phone       string
	bonusPoints int
}

// --- Factory ---

func NewUser(email, password string, role Role) (*User, error) {
	if email == "" || !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}
	if len(password) < 6 {
		return nil, ErrWeakPassword
	}

	id, _ := uuid.NewV7()
	u := &User{
		id:        id.String(),
		email:     email,
		role:      role,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	if err := u.SetPassword(password); err != nil {
		return nil, err
	}

	return u, nil
}

// --- Behavior ---

func (u *User) SetPassword(plainPassword string) error {
	if len(plainPassword) < 6 {
		return ErrWeakPassword
	}
	u.passwordHash = "hash_" + plainPassword
	u.updatedAt = time.Now()
	return nil
}

func (u *User) CheckPassword(plainPassword string) bool {
	return u.passwordHash == "hash_"+plainPassword
}

func (u *User) ChangeRole(newRole Role) {
	u.role = newRole
	u.updatedAt = time.Now()
}

func (u *User) UpdateProfile(name, phone string) {
	u.name = name
	u.phone = phone
	u.updatedAt = time.Now()
}

func (u *User) AccrueBonuses(amount int) {
	if amount > 0 {
		u.bonusPoints += amount
	}
}

func (u *User) SpendBonuses(amount int) error {
	if amount <= 0 {
		return nil
	}
	if u.bonusPoints < amount {
		return ErrInsufficientBonus
	}
	u.bonusPoints -= amount
	return nil
}

func (u *User) ID() string             { return u.id }
func (u *User) Email() string          { return u.email }
func (u *User) Role() Role             { return u.role }
func (u *User) HashedPassword() string { return u.passwordHash }
func (u *User) IsClient() bool         { return u.isClient }
func (u *User) BonusPoints() int       { return u.bonusPoints }
func (u *User) Name() string           { return u.name }
func (u *User) Phone() string          { return u.phone }

type UserRepository interface {
	Save(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}