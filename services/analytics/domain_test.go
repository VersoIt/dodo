package analytics

import (
	"testing"
	"github.com/shopspring/decimal"
)

func TestManagerKPI_CalculateKPIPercent(t *testing.T) {
	kpi := NewManagerKPI("m1", decimal.NewFromInt(1000))

	kpi.AddRevenue(decimal.NewFromInt(500))
	if pct := kpi.CalculateKPIPercent(); !pct.Equal(decimal.NewFromInt(50)) {
		t.Errorf("expected 50%%, got %v%%", pct)
	}

	kpi.AddRevenue(decimal.NewFromInt(500))
	if pct := kpi.CalculateKPIPercent(); !pct.Equal(decimal.NewFromInt(100)) {
		t.Errorf("expected 100%%, got %v%%", pct)
	}

	if !kpi.HasBonus() {
		t.Errorf("manager should have bonus at 100%% plan")
	}
}

func TestManagerKPI_NoPlan(t *testing.T) {
	kpi := NewManagerKPI("m1", decimal.Zero)
	if pct := kpi.CalculateKPIPercent(); !pct.Equal(decimal.NewFromInt(100)) {
		t.Errorf("expected 100%% when no plan, got %v", pct)
	}
}
