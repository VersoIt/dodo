package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/shopspring/decimal"
	"github.com/versoit/diploma/services/analytics"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type AnalyticsUseCase struct {
	repo analytics.AnalyticsRepository
	log  *slog.Logger
}

func NewAnalyticsUseCase(repo analytics.AnalyticsRepository, log *slog.Logger) *AnalyticsUseCase {
	return &AnalyticsUseCase{
		repo: repo,
		log:  log,
	}
}

func (uc *AnalyticsUseCase) RecordSale(ctx context.Context, managerID string, amount decimal.Decimal) error {
	if managerID == "" {
		return fmt.Errorf("%w: manager ID is required", ErrInvalidInput)
	}
	if !amount.IsPositive() {
		return fmt.Errorf("%w: sale amount must be positive", ErrInvalidInput)
	}

	uc.log.Info("recording sale", slog.String("manager_id", managerID), slog.String("amount", amount.String()))

	kpi, err := uc.repo.GetKPI(ctx, managerID)
	if err != nil {
		if errors.Is(err, analytics.ErrKPINotFound) {
			uc.log.Info("creating new kpi for manager", slog.String("manager_id", managerID))
			kpi = analytics.NewManagerKPI(managerID, decimal.NewFromInt(100000))
		} else {
			return fmt.Errorf("failed to retrieve kpi for manager %s: %w", managerID, err)
		}
	}

	kpi.AddRevenue(amount)

	if err := uc.repo.SaveKPI(ctx, kpi); err != nil {
		return fmt.Errorf("failed to update analytics data for manager %s: %w", managerID, err)
	}

	return nil
}

func (uc *AnalyticsUseCase) GetManagerPerformance(ctx context.Context, managerID string) (*analytics.ManagerKPI, error) {
	if managerID == "" {
		return nil, fmt.Errorf("%w: manager ID is required", ErrInvalidInput)
	}

	kpi, err := uc.repo.GetKPI(ctx, managerID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve kpi for manager %s: %w", managerID, err)
	}

	return kpi, nil
}
