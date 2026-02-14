package repository

import (
	"context"
	"time"
	"log/slog"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/orders"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type orderRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
	log  *slog.Logger
}

func NewOrderRepository(pool *pgxpool.Pool, log *slog.Logger) orders.OrderRepository {
	return &orderRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:  log,
	}
}

func (r *orderRepo) Save(ctx context.Context, o *orders.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	addr := o.Address()
	sqlStr, args, err := r.sb.Insert("orders").
		Columns("id", "order_number", "customer_id", "status", "delivery_city", "delivery_street", "delivery_house", "delivery_apartment", "delivery_floor", "delivery_comment", "delivery_price", "discount", "promo_code", "final_price").
		Values(o.ID(), o.OrderNumber(), o.CustomerID(), o.Status(), addr.City, addr.Street, addr.House, addr.Apartment, addr.Floor, addr.Comment, o.DeliveryPrice(), o.Discount(), o.PromoCode(), o.FinalPrice()).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, delivery_price = EXCLUDED.delivery_price, discount = EXCLUDED.discount, final_price = EXCLUDED.final_price").
		ToSql()
	if err != nil {
		return err
	}

	r.log.Debug("saving order", slog.String("order_id", o.ID()), slog.String("number", o.OrderNumber()))

	_, err = tx.Exec(ctx, sqlStr, args...)
	if err != nil {
		r.log.Error("failed to save order", slog.Any("error", err), slog.String("order_id", o.ID()))
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM order_items WHERE order_id = $1", o.ID())
	if err != nil {
		return err
	}

	for _, item := range o.Items() {
		itemID := item.ID()
		sqlStr, args, err := r.sb.Insert("order_items").
			Columns("id", "order_id", "product_id", "product_name", "quantity", "base_price", "size_multiplier").
			Values(itemID, o.ID(), item.ProductID(), item.ProductName(), item.Quantity(), item.BasePrice(), item.Size()).
			ToSql()
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, sqlStr, args...)
		if err != nil {
			return err
		}

		for _, t := range item.Toppings() {
			sqlStr, args, err := r.sb.Insert("order_item_toppings").
				Columns("order_item_id", "name", "price").
				Values(itemID, t.Name, t.Price).
				ToSql()
			if err != nil {
				return err
			}
			_, err = tx.Exec(ctx, sqlStr, args...)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *orderRepo) FindByID(ctx context.Context, id string) (*orders.Order, error) {
	sqlStr, args, err := r.sb.Select("id", "order_number", "customer_id", "status", "created_at", "delivery_city", "delivery_street", "delivery_house", "delivery_apartment", "delivery_floor", "delivery_comment", "delivery_price", "discount", "promo_code", "final_price").
		From("orders").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		oid, number, cid, promo, city, street, house, apartment, floor, comment string
		status                                                                  int
		createdAt                                                               time.Time
		delPrice, discount, finalPrice                                          common.Money
	)

	err = r.pool.QueryRow(ctx, sqlStr, args...).Scan(
		&oid, &number, &cid, &status, &createdAt,
		&city, &street, &house, &apartment, &floor, &comment,
		&delPrice, &discount, &promo, &finalPrice,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, orders.ErrOrderNotFound
		}
		return nil, err
	}

	sqlStr, args, err = r.sb.Select("id", "product_id", "product_name", "quantity", "base_price", "size_multiplier").
		From("order_items").
		Where(squirrel.Eq{"order_id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*orders.OrderItem
	for rows.Next() {
		var (
			iid, pid, name string
			qty            int
			base           common.Money
			size           float64
		)
		if err := rows.Scan(&iid, &pid, &name, &qty, &base, &size); err != nil {
			return nil, err
		}

		tsql, targs, _ := r.sb.Select("name", "price").
			From("order_item_toppings").
			Where(squirrel.Eq{"order_item_id": iid}).
			ToSql()
		
		trows, err := r.pool.Query(ctx, tsql, targs...)
		if err != nil {
			return nil, err
		}
		
		var toppings []orders.Topping
		for trows.Next() {
			var tn string
			var tp common.Money
			if err := trows.Scan(&tn, &tp); err != nil {
				trows.Close()
				return nil, err
			}
			toppings = append(toppings, orders.Topping{Name: tn, Price: tp})
		}
		trows.Close()

		items = append(items, orders.ReconstructOrderItem(iid, pid, name, qty, base, size, toppings))
	}

	addr := orders.DeliveryAddress{
		City: city, Street: street, House: house, Apartment: apartment, Floor: floor, Comment: comment,
	}

	return orders.ReconstructOrder(oid, number, cid, orders.OrderStatus(status), createdAt, addr, delPrice, discount, promo, finalPrice, items), nil
}
