package usecase

import (
	"context"
	"fmt"

	"github.com/versoit/diploma/services/analytics"
)

type AnalyticsUseCase struct {
	repo analytics.AnalyticsRepository
}

func NewAnalyticsUseCase(repo analytics.AnalyticsRepository) *AnalyticsUseCase {
	return &AnalyticsUseCase{repo: repo}
}

// RecordSale добавляет выручку к KPI менеджера.
func (uc *AnalyticsUseCase) RecordSale(ctx context.Context, managerID string, amount float64) error {
	kpi, err := uc.repo.GetKPI(managerID)
	if err != nil {
		// Если KPI на сегодня еще нет, создаем (в реальности план берется из БД)
		kpi = analytics.NewManagerKPI(managerID, 100000) // Пример плана
	}

	kpi.AddRevenue(amount)

	if err := uc.repo.SaveKPI(kpi); err != nil {
		return fmt.Errorf("failed to save kpi: %w", err)
	}

	return nil
}

// GetManagerPerformance возвращает текущие показатели выполнения плана.
func (uc *AnalyticsUseCase) GetManagerPerformance(ctx context.Context, managerID string) (*analytics.ManagerKPI, error) {
	return uc.repo.GetKPI(managerID)
}
