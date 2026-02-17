package domain

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	email, _ := NewEmail("test@example.com")
	pass, _ := NewPassword("password123")
	phone, _ := NewPhone("+79991234567")
	
	u, err := NewUser(email, pass, "Test User", RoleClient, phone)
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
	email, _ := NewEmail("test@example.com")
	pass, _ := NewPassword("password")
	phone, _ := NewPhone("")
	u, _ := NewUser(email, pass, "Test User", RoleClient, phone)

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
