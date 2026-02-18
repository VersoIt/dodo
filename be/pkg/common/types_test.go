package common

import (
	"testing"
)

func TestMoney_Round(t *testing.T) {
	tests := []struct {
		name     string
		input    Money
		expected Money
	}{
		{"No rounding needed", NewMoney(10.50), NewMoney(10.50)},
		{"Round down", NewMoney(10.504), NewMoney(10.50)},
		{"Round up", NewMoney(10.506), NewMoney(10.51)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.Round(2); !got.Equal(tt.expected) {
				t.Errorf("Money.Round(2) = %v, want %v", got, tt.expected)
			}
		})
	}
}
