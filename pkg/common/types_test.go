package common

import (
	"testing"
)

func TestMoney_Round(t *testing.T) {
	tests := []struct {
		name     string
		input    Money
		expected float64
	}{
		{"No rounding needed", 10.50, 10.50},
		{"Round down", 10.504, 10.50},
		{"Round up", 10.506, 10.51}, // Note: current implementation uses int casting which truncates
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.Round(); got != tt.expected {
				t.Errorf("Money.Round() = %v, want %v", got, tt.expected)
			}
		})
	}
}
