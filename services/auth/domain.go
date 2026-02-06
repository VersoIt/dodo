package auth

import (
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

// String возвращает строковое представление роли.
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

// User - Агрегат пользователя.
type User struct {
	id           string
	email        string
	passwordHash string
	role         Role
	createdAt    time.Time
	updatedAt    time.Time

	// Профиль клиента
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

	u := &User{
		id:        uuid.New().String(),
		email:     email,
		role:      role,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	// Хеш будет установлен здесь
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
	// TODO: Реализовать здесь реальное хеширование (bcrypt) в Infrastructure слое,
	// но домен может принимать уже хэш или использовать интерфейс PasswordHasher.
	// Для упрощения пока имитируем:
	u.passwordHash = "hash_" + plainPassword
	u.updatedAt = time.Now()
	return nil
}

func (u *User) CheckPassword(plainPassword string) bool {
	// TODO: compare hash
	return u.passwordHash == "hash_"+plainPassword
}

func (u *User) ChangeRole(newRole Role) {
	u.role = newRole
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

// Getters

func (u *User) ID() string { return u.id }

func (u *User) Email() string { return u.email }

func (u *User) Role() Role { return u.role }

func (u *User) HashedPassword() string { return u.passwordHash }

func (u *User) IsClient() bool { return u.isClient }

func (u *User) BonusPoints() int { return u.bonusPoints }



// --- Repository ---

type UserRepository interface {
	Save(u *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
}
