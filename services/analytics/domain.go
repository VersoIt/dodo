package analytics

import (
	"time"
)

// --- Aggregates ---

// SalesReport - Отчет. Value Object или Entity (зависит от использования).
// Оставим простым DTO, так как это результат вычислений, а не машина состояний.
type SalesReport struct {
	PeriodStart  time.Time
	PeriodEnd    time.Time
	TotalRevenue float64
	TotalOrders  int
	AverageCheck float64
}

// ManagerKPI - Агрегат для расчета премии.
type ManagerKPI struct {
	managerID   string
	shiftDate   time.Time
	planRevenue float64 // Plan_j
	factRevenue float64 // Revenue_fact
}

// --- Factory ---

func NewManagerKPI(managerID string, plan float64) *ManagerKPI {
	return &ManagerKPI{
		managerID:   managerID,
		shiftDate:   time.Now(), // или конкретная дата
		planRevenue: plan,
		factRevenue: 0,
	}
}

// --- Behavior ---

// AddRevenue добавляет выручку к факту.
// Формула (5): Revenue_fact = SUM(Fact_order)
func (k *ManagerKPI) AddRevenue(amount float64) {
	if amount > 0 {
		k.factRevenue += amount
	}
}

// CalculateKPIPercent возвращает процент выполнения плана.
// Формула (6): KPI = (Revenue_fact / Plan_j) * 100%
func (k *ManagerKPI) CalculateKPIPercent() float64 {
	if k.planRevenue == 0 {
		return 100 // Или 0, зависит от логики. Если плана нет - молодцы?
	}
	return (k.factRevenue / k.planRevenue) * 100
}

// HasBonus возвращает true, если бонус положен.
// Формула (7): delta = 1 если KPI >= 100%
func (k *ManagerKPI) HasBonus() bool {
	return k.CalculateKPIPercent() >= 100
}

// Getters
func (k *ManagerKPI) ManagerID() string { return k.managerID }
func (k *ManagerKPI) ShiftDate() time.Time { return k.shiftDate }
func (k *ManagerKPI) Plan() float64 { return k.planRevenue }
func (k *ManagerKPI) Fact() float64 { return k.factRevenue }

// --- Repository ---

type AnalyticsRepository interface {
	SaveKPI(k *ManagerKPI) error
	GetKPI(managerID string) (*ManagerKPI, error)
}
