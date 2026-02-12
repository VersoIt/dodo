package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/analytics"
)

type analyticsRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewAnalyticsRepository(pool *pgxpool.Pool) analytics.AnalyticsRepository {
	return &analyticsRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *analyticsRepo) SaveKPI(ctx context.Context, k *analytics.ManagerKPI) error {
	sql, args, err := r.sb.Insert("manager_kpis").
		Columns("manager_id", "shift_date", "plan_revenue", "fact_revenue").
		Values(k.ManagerID(), k.ShiftDate(), k.Plan(), k.Fact()).
		Suffix("ON CONFLICT (manager_id, shift_date) DO UPDATE SET plan_revenue = EXCLUDED.plan_revenue, fact_revenue = EXCLUDED.fact_revenue").
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *analyticsRepo) GetKPI(ctx context.Context, managerID string) (*analytics.ManagerKPI, error) {
	// Simple query for the latest KPI
	sql, args, err := r.sb.Select("manager_id", "shift_date", "plan_revenue", "fact_revenue").
		From("manager_kpis").
		Where(squirrel.Eq{"manager_id": managerID}).
		OrderBy("shift_date DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		mid  string
		date time.Time
		plan common.Money
		fact common.Money
	)

	err = r.pool.QueryRow(ctx, sql, args...).Scan(&mid, &date, &plan, &fact)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, analytics.ErrKPINotFound
		}
		return nil, err
	}

	return analytics.ReconstructKPI(mid, plan.InexactFloat64(), fact.InexactFloat64()), nil
}
