package repository

import (
	"context"
	"time"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/analytics"
)

type analyticsRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
	log  *slog.Logger
}

func NewAnalyticsRepository(pool *pgxpool.Pool, log *slog.Logger) analytics.AnalyticsRepository {
	return &analyticsRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:  log,
	}
}

func (r *analyticsRepo) SaveKPI(ctx context.Context, k *analytics.ManagerKPI) error {
	sqlStr, args, err := r.sb.Insert("manager_kpis").
		Columns("manager_id", "shift_date", "plan_revenue", "fact_revenue").
		Values(k.ManagerID(), k.ShiftDate(), k.Plan(), k.Fact()).
		Suffix("ON CONFLICT (manager_id, shift_date) DO UPDATE SET plan_revenue = EXCLUDED.plan_revenue, fact_revenue = EXCLUDED.fact_revenue").
		ToSql()
	if err != nil {
		return err
	}

	r.log.Debug("saving manager kpi", slog.String("manager_id", k.ManagerID()), slog.Float64("fact", k.Fact().InexactFloat64()))

	_, err = r.pool.Exec(ctx, sqlStr, args...)
	if err != nil {
		r.log.Error("failed to save kpi", slog.Any("error", err), slog.String("manager_id", k.ManagerID()))
	}
	return err
}

func (r *analyticsRepo) GetKPI(ctx context.Context, managerID string) (*analytics.ManagerKPI, error) {
	sqlStr, args, err := r.sb.Select("manager_id", "shift_date", "plan_revenue", "fact_revenue").
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

	err = r.pool.QueryRow(ctx, sqlStr, args...).Scan(&mid, &date, &plan, &fact)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, analytics.ErrKPINotFound
		}
		return nil, err
	}

	return analytics.ReconstructKPI(mid, plan.InexactFloat64(), fact.InexactFloat64()), nil
}
