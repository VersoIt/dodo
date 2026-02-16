package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/analytics/internal/domain"
)

type analyticsRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	sb     squirrel.StatementBuilderType
	log    *slog.Logger
}

func NewAnalyticsRepository(pool *pgxpool.Pool, log *slog.Logger) domain.AnalyticsRepository {
	return &analyticsRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		sb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:    log,
	}
}

func (r *analyticsRepo) SaveKPI(ctx context.Context, k *domain.ManagerKPI) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Insert("manager_kpis").
		Columns("manager_id", "shift_date", "plan_revenue", "fact_revenue").
		Values(k.ManagerID(), k.ShiftDate(), k.Plan(), k.Fact()).
		Suffix("ON CONFLICT (manager_id, shift_date) DO UPDATE SET plan_revenue = EXCLUDED.plan_revenue, fact_revenue = EXCLUDED.fact_revenue").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	r.log.Debug("saving manager kpi", slog.String("manager_id", k.ManagerID()), slog.Float64("fact", k.Fact().InexactFloat64()))

	if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("exec save kpi: %w", err)
	}
	return nil
}

func (r *analyticsRepo) GetKPI(ctx context.Context, managerID string) (*domain.ManagerKPI, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("manager_id", "shift_date", "plan_revenue", "fact_revenue").
		From("manager_kpis").
		Where(squirrel.Eq{"manager_id": managerID}).
		OrderBy("shift_date DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var (
		mid  string
		date time.Time
		plan common.Money
		fact common.Money
	)

	err = db.QueryRow(ctx, sqlStr, args...).Scan(&mid, &date, &plan, &fact)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrKPINotFound
		}
		return nil, fmt.Errorf("query kpi: %w", err)
	}

	return domain.ReconstructKPI(mid, plan.InexactFloat64(), fact.InexactFloat64()), nil
}
