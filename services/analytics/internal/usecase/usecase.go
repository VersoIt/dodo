package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/shopspring/decimal"
	"github.com/versoit/diploma/services/analytics/internal/domain"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type AnalyticsUseCase struct {
	repo domain.AnalyticsRepository
	tm   trm.Manager
	log  *slog.Logger
}

func NewAnalyticsUseCase(repo domain.AnalyticsRepository, tm trm.Manager, log *slog.Logger) *AnalyticsUseCase {
	return &AnalyticsUseCase{
		repo: repo,
		tm:   tm,
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

	return uc.tm.Do(ctx, func(ctx context.Context) error {
		kpi, err := uc.repo.GetKPI(ctx, managerID)
		if err != nil {
			if errors.Is(err, domain.ErrKPINotFound) {
				uc.log.Info("creating new kpi for manager", slog.String("manager_id", managerID))
				kpi = domain.NewManagerKPI(managerID, decimal.NewFromInt(100000))
			} else {
				return fmt.Errorf("retrieve kpi for manager %s: %w", managerID, err)
			}
		}

		kpi.AddRevenue(amount)

		if err := uc.repo.SaveKPI(ctx, kpi); err != nil {
			return fmt.Errorf("update analytics data for manager %s: %w", managerID, err)
		}
		return nil
	})
}

func (uc *AnalyticsUseCase) GetManagerPerformance(ctx context.Context, managerID string) (*domain.ManagerKPI, error) {
	if managerID == "" {
		return nil, fmt.Errorf("%w: manager ID is required", ErrInvalidInput)
	}

	kpi, err := uc.repo.GetKPI(ctx, managerID)
	if err != nil {
		if errors.Is(err, domain.ErrKPINotFound) {
			return nil, domain.ErrKPINotFound
		}
		return nil, fmt.Errorf("retrieve kpi for manager %s: %w", managerID, err)
	}

	return kpi, nil
}
