package analytics

import (
	"testing"
)

func TestManagerKPI_CalculateKPIPercent(t *testing.T) {
	kpi := NewManagerKPI("m1", 1000)
	
	kpi.AddRevenue(500)
	if pct := kpi.CalculateKPIPercent(); pct != 50 {
		t.Errorf("expected 50%%, got %v%%", pct)
	}

	kpi.AddRevenue(500)
	if pct := kpi.CalculateKPIPercent(); pct != 100 {
		t.Errorf("expected 100%%, got %v%%", pct)
	}
	
	if !kpi.HasBonus() {
		t.Errorf("manager should have bonus at 100%% plan")
	}
}

func TestManagerKPI_NoPlan(t *testing.T) {
	kpi := NewManagerKPI("m1", 0)
	if pct := kpi.CalculateKPIPercent(); pct != 100 {
		t.Errorf("expected 100%% when no plan, got %v", pct)
	}
}
