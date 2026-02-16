package domain

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// --- Email ---

type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func NewEmail(v string) (Email, error) {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return Email{}, fmt.Errorf("email cannot be empty")
	}
	if !emailRegex.MatchString(v) {
		return Email{}, fmt.Errorf("invalid email format")
	}
	return Email{value: v}, nil
}

func (e Email) String() string {
	return e.value
}

// --- Password ---

type Password struct {
	value string
}

func NewPassword(v string) (Password, error) {
	if utf8.RuneCountInString(v) < 6 {
		return Password{}, fmt.Errorf("password must be at least 6 characters")
	}
	return Password{value: v}, nil
}

func (p Password) String() string {
	return p.value
}

// --- Phone ---

type Phone struct {
	value string
}

func NewPhone(v string) (Phone, error) {
	v = strings.TrimSpace(v)
	// Allow empty phone for now if it's optional, or enforce rules
	if v != "" {
		// Basic check: must contain digits
		digits := 0
		for _, r := range v {
			if r >= '0' && r <= '9' {
				digits++
			}
		}
		if digits < 7 {
			return Phone{}, fmt.Errorf("phone number too short")
		}
	}
	return Phone{value: v}, nil
}

func (p Phone) String() string {
	return p.value
}
