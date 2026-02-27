package common

import "github.com/shopspring/decimal"

// Money - денежный тип для всей системы.
// Используем decimal.Decimal для предотвращения ошибок округления floating-point.
type Money = decimal.Decimal

// NewMoney - хелпер для создания Money из float64.
func NewMoney(v float64) Money {
	return decimal.NewFromFloat(v)
}

// ZeroMoney - нулевое значение.
func ZeroMoney() Money {
	return decimal.NewFromInt(0)
}
