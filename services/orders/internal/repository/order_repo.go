package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/orders"
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

// --- Order Persistence ---

func (r *orderRepo) Save(ctx context.Context, o *orders.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil { return err }
	defer func() { _ = tx.Rollback(ctx) }()

	addr := o.Address()
	var chefID, courierID *string
	if o.ChefID() != "" { cid := o.ChefID(); chefID = &cid }
	if o.CourierID() != "" { coid := o.CourierID(); courierID = &coid }

	sqlStr, args, err := r.sb.Insert("orders").
		Columns("id", "order_number", "customer_id", "status", "delivery_city", "delivery_street", "delivery_house", "delivery_apartment", "delivery_floor", "delivery_entrance", "delivery_comment", "delivery_price", "discount", "promo_code", "final_price", "chef_id", "courier_id").
		Values(o.ID(), o.OrderNumber(), o.CustomerID(), o.Status(), addr.City, addr.Street, addr.House, addr.Apartment, addr.Floor, addr.Entrance, addr.Comment, o.DeliveryPrice(), o.Discount(), o.PromoCode(), o.FinalPrice(), chefID, courierID).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, delivery_price = EXCLUDED.delivery_price, discount = EXCLUDED.discount, final_price = EXCLUDED.final_price, chef_id = EXCLUDED.chef_id, courier_id = EXCLUDED.courier_id").
		ToSql()
	if err != nil { return err }

	if _, err = tx.Exec(ctx, sqlStr, args...); err != nil { return err }
	if _, err = tx.Exec(ctx, "DELETE FROM order_items WHERE order_id = $1", o.ID()); err != nil { return err }

	for _, item := range o.Items() {
		itemID := item.ID()
		sqlStr, args, err := r.sb.Insert("order_items").
			Columns("id", "order_id", "product_id", "product_name", "quantity", "base_price", "size_multiplier").
			Values(itemID, o.ID(), item.ProductID(), item.ProductName(), item.Quantity(), item.BasePrice(), item.Size()).
			ToSql()
		if err != nil { return err }
		if _, err = tx.Exec(ctx, sqlStr, args...); err != nil { return err }

		for _, t := range item.Toppings() {
			sqlStr, args, err := r.sb.Insert("order_item_toppings").
				Columns("order_item_id", "name", "price").
				Values(itemID, t.Name, t.Price).
				ToSql()
			if err != nil { return err }
			if _, err = tx.Exec(ctx, sqlStr, args...); err != nil { return err }
		}
	}
	return tx.Commit(ctx)
}

func (r *orderRepo) FindByID(ctx context.Context, id string) (*orders.Order, error) {
	sqlStr, args, err := r.sb.Select("id", "order_number", "customer_id", "status", "created_at", "delivery_city", "delivery_street", "delivery_house", "delivery_apartment", "delivery_floor", "delivery_entrance", "delivery_comment", "delivery_price", "discount", "promo_code", "final_price", "chef_id", "courier_id").
		From("orders").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil { return nil, err }

	var (
		oid, number, cid, promo, city, street, house, apartment, floor, entrance, comment string
		status int
		createdAt time.Time
		delPrice, discount, finalPrice common.Money
		chefID, courierID sql.NullString
	)

	err = r.pool.QueryRow(ctx, sqlStr, args...).Scan(
		&oid, &number, &cid, &status, &createdAt, &city, &street, &house, &apartment, &floor, &entrance, &comment,
		&delPrice, &discount, &promo, &finalPrice, &chefID, &courierID,
	)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, orders.ErrOrderNotFound }
		return nil, err
	}

	rows, _ := r.pool.Query(ctx, "SELECT id, product_id, product_name, quantity, base_price, size_multiplier FROM order_items WHERE order_id = $1", oid)
	defer rows.Close()
	var items []*orders.OrderItem
	for rows.Next() {
		var (iid, pid, name string; qty int; base common.Money; size float64)
		if err := rows.Scan(&iid, &pid, &name, &qty, &base, &size); err == nil {
			var toppings []orders.Topping
			trows, _ := r.pool.Query(ctx, "SELECT name, price FROM order_item_toppings WHERE order_item_id = $1", iid)
			for trows.Next() {
				var tn string; var tp common.Money
				if err := trows.Scan(&tn, &tp); err == nil { toppings = append(toppings, orders.Topping{Name: tn, Price: tp}) }
			}
			trows.Close()
			items = append(items, orders.ReconstructOrderItem(iid, pid, name, qty, base, size, toppings))
		}
	}
	addr := orders.DeliveryAddress{City: city, Street: street, House: house, Apartment: apartment, Floor: floor, Entrance: entrance, Comment: comment}
	return orders.ReconstructOrder(oid, number, cid, orders.OrderStatus(status), createdAt, addr, delPrice, discount, promo, finalPrice, items, chefID.String, courierID.String), nil
}

func (r *orderRepo) FindByCustomerID(ctx context.Context, customerID string) ([]*orders.Order, error) {
	rows, _ := r.pool.Query(ctx, "SELECT id FROM orders WHERE customer_id = $1 ORDER BY created_at DESC", customerID)
	defer rows.Close()
	var result []*orders.Order
	for rows.Next() {
		var id string; _ = rows.Scan(&id)
		if o, err := r.FindByID(ctx, id); err == nil { result = append(result, o) }
	}
	return result, nil
}

func (r *orderRepo) FindAll(ctx context.Context) ([]*orders.Order, error) {
	rows, _ := r.pool.Query(ctx, "SELECT id FROM orders ORDER BY created_at DESC")
	defer rows.Close()
	var result []*orders.Order
	for rows.Next() {
		var id string; _ = rows.Scan(&id)
		if o, err := r.FindByID(ctx, id); err == nil { result = append(result, o) }
	}
	return result, nil
}

// --- Promo Codes ---

func (r *orderRepo) SavePromo(ctx context.Context, p *orders.PromoCode) error {
	r.log.Info("Saving promo code", slog.String("code", p.Code()))
	sqlStr, args, err := r.sb.Insert("promo_codes").
		Columns("id", "code", "discount_type", "discount_amount", "is_active").
		Values(p.ID(), p.Code(), p.DiscountType(), p.DiscountAmount(), p.IsActive()).
		Suffix("ON CONFLICT (code) DO UPDATE SET discount_amount = EXCLUDED.discount_amount, is_active = EXCLUDED.is_active").
		ToSql()
	if err != nil { return err }
	_, err = r.pool.Exec(ctx, sqlStr, args...)
	return err
}

func (r *orderRepo) FindPromoByCode(ctx context.Context, code string) (*orders.PromoCode, error) {
	var (id, c, dType string; amount common.Money; active bool; expires sql.NullTime)
	err := r.pool.QueryRow(ctx, "SELECT id, code, discount_type, discount_amount, is_active, expires_at FROM promo_codes WHERE code = $1", code).
		Scan(&id, &c, &dType, &amount, &active, &expires)
	if err != nil { return nil, err }
	return orders.NewPromoCode(id, c, dType, amount, active, expires.Time), nil
}

func (r *orderRepo) ListPromos(ctx context.Context) ([]*orders.PromoCode, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, code, discount_type, discount_amount, is_active, expires_at FROM promo_codes ORDER BY created_at DESC")
	if err != nil { return nil, err }
	defer rows.Close()
	var res []*orders.PromoCode
	for rows.Next() {
		var (id, c, dType string; amount common.Money; active bool; expires sql.NullTime)
		if err := rows.Scan(&id, &c, &dType, &amount, &active, &expires); err != nil {
			r.log.Error("Failed to scan promo", slog.Any("error", err))
			continue
		}
		res = append(res, orders.NewPromoCode(id, c, dType, amount, active, expires.Time))
	}
	return res, nil
}

func (r *orderRepo) DeletePromo(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM promo_codes WHERE id = $1", id)
	return err
}

// --- Analytics ---

func (r *orderRepo) GetKPIs(ctx context.Context) (*orders.OrderStats, error) {
	var stats orders.OrderStats
	err := r.pool.QueryRow(ctx, `
		SELECT 
			COALESCE(SUM(final_price), 0) as total_revenue,
			COUNT(*) as orders_count,
			COALESCE(AVG(final_price), 0) as avg_check
		FROM orders 
		WHERE status != 6
	`).Scan(&stats.TotalRevenue, &stats.OrdersCount, &stats.AvgCheck)
	return &stats, err
}

func (r *orderRepo) GetTopProducts(ctx context.Context, limit int) ([]orders.ProductStat, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT 
			product_name, 
			SUM(quantity) as count, 
			SUM(base_price * quantity) as revenue
		FROM order_items
		GROUP BY product_name
		ORDER BY count DESC
		LIMIT $1
	`, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var stats []orders.ProductStat
	for rows.Next() {
		var s orders.ProductStat
		if err := rows.Scan(&s.Name, &s.Count, &s.Revenue); err == nil { stats = append(stats, s) }
	}
	return stats, nil
}
