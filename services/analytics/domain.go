package analytics

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type SalesReport struct {
	PeriodStart  time.Time
	PeriodEnd    time.Time
	TotalRevenue decimal.Decimal
	TotalOrders  int
	AverageCheck decimal.Decimal
}

type ManagerKPI struct {
	managerID   string
	shiftDate   time.Time
	planRevenue decimal.Decimal
	factRevenue decimal.Decimal
}

func NewManagerKPI(managerID string, plan decimal.Decimal) *ManagerKPI {
	return &ManagerKPI{
		managerID:   managerID,
		shiftDate:   time.Now(),
		planRevenue: plan,
		factRevenue: decimal.Zero,
	}
}

func (k *ManagerKPI) AddRevenue(amount decimal.Decimal) {
	if amount.IsPositive() {
		k.factRevenue = k.factRevenue.Add(amount)
	}
}

func (k *ManagerKPI) CalculateKPIPercent() decimal.Decimal {
	if k.planRevenue.IsZero() {
		return decimal.NewFromInt(100)
	}
	return k.factRevenue.Div(k.planRevenue).Mul(decimal.NewFromInt(100))
}

func (k *ManagerKPI) HasBonus() bool {
	return k.CalculateKPIPercent().GreaterThanOrEqual(decimal.NewFromInt(100))
}

func (k *ManagerKPI) ManagerID() string         { return k.managerID }
func (k *ManagerKPI) ShiftDate() time.Time      { return k.shiftDate }
func (k *ManagerKPI) Plan() decimal.Decimal     { return k.planRevenue }
func (k *ManagerKPI) Fact() decimal.Decimal     { return k.factRevenue }

type AnalyticsRepository interface {
	SaveKPI(ctx context.Context, k *ManagerKPI) error
	GetKPI(ctx context.Context, managerID string) (*ManagerKPI, error)
}
