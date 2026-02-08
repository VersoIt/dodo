package repository

import (
	"context"

	"github.com/versoit/diploma/services/logistics"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type deliveryRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewDeliveryRepository(pool *pgxpool.Pool) logistics.DeliveryRepository {
	return &deliveryRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *deliveryRepo) Save(ctx context.Context, d *logistics.Delivery) error {
	sql, args, err := r.sb.Insert("deliveries").
		Columns("id", "order_id", "courier_id", "status").
		Values(d.ID(), d.OrderID(), d.CourierID(), d.Status()).
		Suffix("ON CONFLICT (id) DO UPDATE SET courier_id = EXCLUDED.courier_id, status = EXCLUDED.status, updated_at = NOW()").
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *deliveryRepo) FindByOrderID(ctx context.Context, orderID string) (*logistics.Delivery, error) {
	return nil, nil
}

type courierRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewCourierRepository(pool *pgxpool.Pool) logistics.CourierRepository {
	return &courierRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *courierRepo) Save(ctx context.Context, c *logistics.Courier) error {
	sql, args, err := r.sb.Insert("couriers").
		Columns("id", "name", "phone", "status").
		Values(c.ID(), c.Name(), c.Phone(), c.Status()).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, updated_at = NOW()").
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *courierRepo) FindByID(ctx context.Context, id string) (*logistics.Courier, error) {
	return nil, nil
}

func (r *courierRepo) FindAvailable(ctx context.Context) ([]*logistics.Courier, error) {
	return nil, nil
}

func (r *courierRepo) UpdateLocation(ctx context.Context, id string, lat, lng float64) error {
	return nil
}

