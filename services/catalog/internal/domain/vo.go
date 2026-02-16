package domain

import (
	"fmt"
	"strings"

	"github.com/versoit/diploma/pkg/common"
)

// --- ProductName ---

type ProductName struct {
	value string
}

func NewProductName(v string) (ProductName, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return ProductName{}, fmt.Errorf("product name cannot be empty")
	}
	return ProductName{value: v}, nil
}

func (n ProductName) String() string {
	return n.value
}

// --- Price ---

type Price struct {
	amount common.Money
}

func NewPrice(amount float64) (Price, error) {
	money := common.NewMoney(amount)
	if money.IsNegative() {
		return Price{}, fmt.Errorf("price cannot be negative")
	}
	return Price{amount: money}, nil
}

// FromMoney allows creating Price from existing Money object, validating it
func NewPriceFromMoney(m common.Money) (Price, error) {
	if m.IsNegative() {
		return Price{}, fmt.Errorf("price cannot be negative")
	}
	return Price{amount: m}, nil
}

func (p Price) Money() common.Money {
	return p.amount
}
