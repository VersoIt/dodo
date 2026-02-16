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
	"github.com/versoit/diploma/services/logistics/internal/domain"
)

type deliveryRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	sb     squirrel.StatementBuilderType
	log    *slog.Logger
}

func NewDeliveryRepository(pool *pgxpool.Pool, log *slog.Logger) domain.DeliveryRepository {
	return &deliveryRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		sb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:    log,
	}
}

func (r *deliveryRepo) Save(ctx context.Context, d *domain.Delivery) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	lat, lng := d.Location()
	sqlStr, args, err := r.sb.Insert("deliveries").
		Columns("order_id", "courier_id", "status", "current_lat", "current_lng").
		Values(d.OrderID(), d.CourierID(), d.Status(), lat, lng).
		Suffix("ON CONFLICT (order_id) DO UPDATE SET courier_id = EXCLUDED.courier_id, status = EXCLUDED.status, current_lat = EXCLUDED.current_lat, current_lng = EXCLUDED.current_lng, updated_at = NOW()").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	r.log.Debug("saving delivery", slog.String("order_id", d.OrderID()), slog.Int("status", int(d.Status())))

	if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("exec save delivery: %w", err)
	}
	return nil
}

func (r *deliveryRepo) FindByOrderID(ctx context.Context, orderID string) (*domain.Delivery, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("order_id", "courier_id", "status", "created_at", "pickup_time", "delivery_time", "current_lat", "current_lng").
		From("deliveries").
		Where(squirrel.Eq{"order_id": orderID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var (
		oid      string
		cid      *string
		status   int
		ca       time.Time
		pt, dt   *time.Time
		lat, lng float64
	)

	err = db.QueryRow(ctx, sqlStr, args...).Scan(&oid, &cid, &status, &ca, &pt, &dt, &lat, &lng)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("delivery not found")
		}
		return nil, fmt.Errorf("query delivery: %w", err)
	}

	var p, d time.Time
	if pt != nil {
		p = *pt
	}
	if dt != nil {
		d = *dt
	}

	courierID := ""
	if cid != nil {
		courierID = *cid
	}

	return domain.ReconstructDelivery(oid, courierID, domain.DeliveryStatus(status), ca, p, d, lat, lng), nil
}

type courierRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	sb     squirrel.StatementBuilderType
	log    *slog.Logger
}

func NewCourierRepository(pool *pgxpool.Pool, log *slog.Logger) domain.CourierRepository {
	return &courierRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		sb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:    log,
	}
}

func (r *courierRepo) Save(ctx context.Context, c *domain.Courier) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	lat, lng := c.Location()
	sqlStr, args, err := r.sb.Insert("couriers").
		Columns("id", "name", "phone", "status", "current_lat", "current_lng").
		Values(c.ID(), c.Name(), c.Phone(), c.Status(), lat, lng).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, current_lat = EXCLUDED.current_lat, current_lng = EXCLUDED.current_lng, updated_at = NOW()").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	r.log.Debug("saving courier", slog.String("courier_id", c.ID()))

	if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("exec save courier: %w", err)
	}
	return nil
}

func (r *courierRepo) FindByID(ctx context.Context, id string) (*domain.Courier, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("id", "name", "phone", "status", "current_lat", "current_lng").
		From("couriers").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var (
		cid, name, phone string
		status           int
		lat, lng         float64
	)

	err = db.QueryRow(ctx, sqlStr, args...).Scan(&cid, &name, &phone, &status, &lat, &lng)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("courier not found")
		}
		return nil, fmt.Errorf("query courier: %w", err)
	}

	return domain.ReconstructCourier(cid, name, phone, domain.CourierStatus(status), lat, lng), nil
}

func (r *courierRepo) FindAvailable(ctx context.Context) ([]*domain.Courier, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("id", "name", "phone", "status", "current_lat", "current_lng").
		From("couriers").
		Where(squirrel.Eq{"status": domain.CourierFree}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query available couriers: %w", err)
	}
	defer rows.Close()

	var couriers []*domain.Courier
	for rows.Next() {
		var (
			cid, name, phone string
			status           int
			lat, lng         float64
		)
		if err := rows.Scan(&cid, &name, &phone, &status, &lat, &lng); err != nil {
			return nil, fmt.Errorf("scan courier: %w", err)
		}
		couriers = append(couriers, domain.ReconstructCourier(cid, name, phone, domain.CourierStatus(status), lat, lng))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return couriers, nil
}

func (r *courierRepo) UpdateLocation(ctx context.Context, id string, lat, lng float64) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Update("couriers").
		Set("current_lat", lat).
		Set("current_lng", lng).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}
	if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("exec update location: %w", err)
	}
	return nil
}
