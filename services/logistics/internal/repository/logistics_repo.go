package repository

import (
	"context"
	"fmt"
	"time"
	"log/slog"

	"github.com/versoit/diploma/services/logistics"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type deliveryRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
	log  *slog.Logger
}

func NewDeliveryRepository(pool *pgxpool.Pool, log *slog.Logger) logistics.DeliveryRepository {
	return &deliveryRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:  log,
	}
}

func (r *deliveryRepo) Save(ctx context.Context, d *logistics.Delivery) error {
	lat, lng := d.Location()
	sqlStr, args, err := r.sb.Insert("deliveries").
		Columns("order_id", "courier_id", "status", "current_lat", "current_lng").
		Values(d.OrderID(), d.CourierID(), d.Status(), lat, lng).
		Suffix("ON CONFLICT (order_id) DO UPDATE SET courier_id = EXCLUDED.courier_id, status = EXCLUDED.status, current_lat = EXCLUDED.current_lat, current_lng = EXCLUDED.current_lng, updated_at = NOW()").
		ToSql()
	if err != nil {
		return err
	}
	
	r.log.Debug("saving delivery", slog.String("order_id", d.OrderID()), slog.Int("status", int(d.Status())))
	
	_, err = r.pool.Exec(ctx, sqlStr, args...)
	if err != nil {
		r.log.Error("failed to save delivery", slog.Any("error", err), slog.String("order_id", d.OrderID()))
	}
	return err
}

func (r *deliveryRepo) FindByOrderID(ctx context.Context, orderID string) (*logistics.Delivery, error) {
	sqlStr, args, err := r.sb.Select("order_id", "courier_id", "status", "created_at", "pickup_time", "delivery_time", "current_lat", "current_lng").
		From("deliveries").
		Where(squirrel.Eq{"order_id": orderID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		oid            string
		cid            *string
		status         int
		ca             time.Time
		pt, dt         *time.Time
		lat, lng       float64
	)

	err = r.pool.QueryRow(ctx, sqlStr, args...).Scan(&oid, &cid, &status, &ca, &pt, &dt, &lat, &lng)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("delivery not found")
		}
		return nil, err
	}

	var p, d time.Time
	if pt != nil { p = *pt }
	if dt != nil { d = *dt }

	courierID := ""
	if cid != nil { courierID = *cid }

	return logistics.ReconstructDelivery(oid, courierID, logistics.DeliveryStatus(status), ca, p, d, lat, lng), nil
}

type courierRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
	log  *slog.Logger
}

func NewCourierRepository(pool *pgxpool.Pool, log *slog.Logger) logistics.CourierRepository {
	return &courierRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:  log,
	}
}

func (r *courierRepo) Save(ctx context.Context, c *logistics.Courier) error {
	lat, lng := c.Location()
	sqlStr, args, err := r.sb.Insert("couriers").
		Columns("id", "name", "phone", "status", "current_lat", "current_lng").
		Values(c.ID(), c.Name(), c.Phone(), c.Status(), lat, lng).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, current_lat = EXCLUDED.current_lat, current_lng = EXCLUDED.current_lng, updated_at = NOW()").
		ToSql()
	if err != nil {
		return err
	}
	
	r.log.Debug("saving courier", slog.String("courier_id", c.ID()))
	
	_, err = r.pool.Exec(ctx, sqlStr, args...)
	return err
}

func (r *courierRepo) FindByID(ctx context.Context, id string) (*logistics.Courier, error) {
	sqlStr, args, err := r.sb.Select("id", "name", "phone", "status", "current_lat", "current_lng").
		From("couriers").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		cid, name, phone string
		status           int
		lat, lng         float64
	)

	err = r.pool.QueryRow(ctx, sqlStr, args...).Scan(&cid, &name, &phone, &status, &lat, &lng)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("courier not found")
		}
		return nil, err
	}

	return logistics.ReconstructCourier(cid, name, phone, logistics.CourierStatus(status), lat, lng), nil
}

func (r *courierRepo) FindAvailable(ctx context.Context) ([]*logistics.Courier, error) {
	sqlStr, args, err := r.sb.Select("id", "name", "phone", "status", "current_lat", "current_lng").
		From("couriers").
		Where(squirrel.Eq{"status": logistics.CourierFree}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var couriers []*logistics.Courier
	for rows.Next() {
		var (
			cid, name, phone string
			status           int
			lat, lng         float64
		)
		if err := rows.Scan(&cid, &name, &phone, &status, &lat, &lng); err != nil {
			return nil, err
		}
		couriers = append(couriers, logistics.ReconstructCourier(cid, name, phone, logistics.CourierStatus(status), lat, lng))
	}
	return couriers, nil
}

func (r *courierRepo) UpdateLocation(ctx context.Context, id string, lat, lng float64) error {
	sqlStr, args, err := r.sb.Update("couriers").
		Set("current_lat", lat).
		Set("current_lng", lng).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, sqlStr, args...)
	return err
}
