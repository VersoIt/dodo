package auth

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("test@example.com", "password123", RoleClient)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if u.Email() != "test@example.com" {
		t.Errorf("email mismatch")
	}

	if !u.CheckPassword("password123") {
		t.Errorf("password check failed")
	}
}

func TestUser_BonusSystem(t *testing.T) {
	u, _ := NewUser("test@example.com", "password", RoleClient)
	
	u.AccrueBonuses(100)
	if u.BonusPoints() != 100 {
		t.Errorf("expected 100 bonuses, got %d", u.BonusPoints())
	}

	err := u.SpendBonuses(40)
	if err != nil || u.BonusPoints() != 60 {
		t.Errorf("failed to spend bonuses")
	}

	err = u.SpendBonuses(100)
	if err != ErrInsufficientBonus {
		t.Errorf("expected ErrInsufficientBonus, got %v", err)
	}
}
