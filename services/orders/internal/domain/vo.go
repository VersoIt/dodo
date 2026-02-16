package domain

import (
	"fmt"
	"strings"

	"github.com/versoit/diploma/pkg/common"
)

// --- Quantity ---

type Quantity struct {
	value int
}

func NewQuantity(v int) (Quantity, error) {
	if v <= 0 {
		return Quantity{}, fmt.Errorf("quantity must be positive")
	}
	return Quantity{value: v}, nil
}

func (q Quantity) Int() int {
	return q.value
}

// --- Address ---

// We reuse DeliveryAddress structure but add a constructor with validation
func NewDeliveryAddress(city, street, house, apartment, floor, entrance, comment string) (DeliveryAddress, error) {
	city = strings.TrimSpace(city)
	street = strings.TrimSpace(street)
	house = strings.TrimSpace(house)

	if city == "" {
		return DeliveryAddress{}, fmt.Errorf("city cannot be empty")
	}
	if street == "" {
		return DeliveryAddress{}, fmt.Errorf("street cannot be empty")
	}
	if house == "" {
		return DeliveryAddress{}, fmt.Errorf("house cannot be empty")
	}

	return DeliveryAddress{
		City:      city,
		Street:    street,
		House:     house,
		Apartment: strings.TrimSpace(apartment),
		Floor:     strings.TrimSpace(floor),
		Entrance:  strings.TrimSpace(entrance),
		Comment:   strings.TrimSpace(comment),
	}, nil
}

// --- Money Wrapper (optional, but good for domain semantics) ---
// For now, we stick to common.Money as it is already a Value Object (decimal.Decimal is immutable)
type Money = common.Money
