package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/be/orders/internal/domain"
)

type orderRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	sb     squirrel.StatementBuilderType
	log    *slog.Logger
}

func NewOrderRepository(pool *pgxpool.Pool, log *slog.Logger) domain.OrderRepository {
	return &orderRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		sb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:    log,
	}
}

// --- Order Persistence ---

func (r *orderRepo) Save(ctx context.Context, o *domain.Order) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	addr := o.Address()
	var chefID, courierID *string
	if o.ChefID() != "" {
		cid := o.ChefID()
		chefID = &cid
	}
	if o.CourierID() != "" {
		coid := o.CourierID()
		courierID = &coid
	}

	sqlStr, args, err := r.sb.Insert("orders").
		Columns("id", "order_number", "customer_id", "status", "delivery_city", "delivery_street", "delivery_house", "delivery_apartment", "delivery_floor", "delivery_entrance", "delivery_comment", "delivery_price", "discount", "promo_code", "final_price", "chef_id", "courier_id").
		Values(o.ID(), o.OrderNumber(), o.CustomerID(), o.Status(), addr.City, addr.Street, addr.House, addr.Apartment, addr.Floor, addr.Entrance, addr.Comment, o.DeliveryPrice(), o.Discount(), o.PromoCode(), o.FinalPrice(), chefID, courierID).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, delivery_price = EXCLUDED.delivery_price, discount = EXCLUDED.discount, final_price = EXCLUDED.final_price, chef_id = EXCLUDED.chef_id, courier_id = EXCLUDED.courier_id").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("exec save order: %w", err)
	}
		if _, err = db.Exec(ctx, "DELETE FROM order_items WHERE order_id = $1", o.ID()); err != nil {
			return fmt.Errorf("delete old items: %w", err)
		}
	
			for _, item := range o.Items() {
	
				itemID := item.ID()
	
		
	
				sqlStr, args, err = r.sb.Insert("order_items").
	
		
				Columns("id", "order_id", "product_id", "product_name", "quantity", "base_price", "size_multiplier").
				Values(itemID, o.ID(), item.ProductID(), item.ProductName(), item.Quantity(), item.BasePrice(), item.Size()).
				ToSql()
		if err != nil {
			return fmt.Errorf("build item query: %w", err)
		}
		if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
			return fmt.Errorf("exec save item: %w", err)
		}

		for _, t := range item.Toppings() {
			sqlStr, args, err = r.sb.Insert("order_item_toppings").
				Columns("order_item_id", "name", "price").
				Values(itemID, t.Name, t.Price).
				ToSql()
			if err != nil {
				return fmt.Errorf("build topping query: %w", err)
			}
			if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
				return fmt.Errorf("exec save topping: %w", err)
			}
		}
	}
	return nil
}

func (r *orderRepo) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("id", "order_number", "customer_id", "status", "created_at", "delivery_city", "delivery_street", "delivery_house", "delivery_apartment", "delivery_floor", "delivery_entrance", "delivery_comment", "delivery_price", "discount", "promo_code", "final_price", "chef_id", "courier_id").
		From("orders").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var (
		oid, number, cid, promo, city, street, house, apartment, floor, entrance, comment string
		status                                                                            int
		createdAt                                                                         time.Time
		delPrice, discount, finalPrice                                                    common.Money
		chefID, courierID                                                                 sql.NullString
	)

	err = db.QueryRow(ctx, sqlStr, args...).Scan(
		&oid, &number, &cid, &status, &createdAt, &city, &street, &house, &apartment, &floor, &entrance, &comment,
		&delPrice, &discount, &promo, &finalPrice, &chefID, &courierID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	// Fetch all items for this order
	itemRows, err := db.Query(ctx, `
		SELECT i.id, i.product_id, i.product_name, i.quantity, i.base_price, i.size_multiplier,
		       t.name, t.price
		FROM order_items i
		LEFT JOIN order_item_toppings t ON i.id = t.order_item_id
		WHERE i.order_id = $1
	`, oid)
	if err != nil {
		return nil, fmt.Errorf("query items and toppings: %w", err)
	}
	defer itemRows.Close()

	itemsMap := make(map[string]*domain.OrderItem)
	var itemIDs []string // to keep order

	for itemRows.Next() {
		var (
			iid, pid, pname string
			qty             int
			base            common.Money
			size            float64
			tname           sql.NullString
			tprice          *common.Money
		)
		if err := itemRows.Scan(&iid, &pid, &pname, &qty, &base, &size, &tname, &tprice); err != nil {
			return nil, fmt.Errorf("scan item/topping: %w", err)
		}

		item, ok := itemsMap[iid]
		if !ok {
			item = domain.ReconstructOrderItem(iid, pid, pname, qty, base, size, nil)
			itemsMap[iid] = item
			itemIDs = append(itemIDs, iid)
		}

		if tname.Valid && tprice != nil {
			item.AddReconstructedTopping(tname.String, *tprice)
		}
	}

	var items []*domain.OrderItem
	for _, id := range itemIDs {
		items = append(items, itemsMap[id])
	}

	addr := domain.DeliveryAddress{City: city, Street: street, House: house, Apartment: apartment, Floor: floor, Entrance: entrance, Comment: comment}
	return domain.ReconstructOrder(oid, number, cid, domain.OrderStatus(status), createdAt, addr, delPrice, discount, promo, finalPrice, items, chefID.String, courierID.String), nil
}

func (r *orderRepo) FindByCustomerID(ctx context.Context, customerID string) ([]*domain.Order, error) {
	return r.FindFiltered(ctx, domain.OrderFilter{CustomerID: &customerID})
}

func (r *orderRepo) FindAll(ctx context.Context) ([]*domain.Order, error) {
	return r.FindFiltered(ctx, domain.OrderFilter{})
}

func (r *orderRepo) FindFiltered(ctx context.Context, filter domain.OrderFilter) ([]*domain.Order, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	builder := r.sb.Select("id", "order_number", "customer_id", "status", "created_at", "delivery_city", "delivery_street", "delivery_house", "delivery_apartment", "delivery_floor", "delivery_entrance", "delivery_comment", "delivery_price", "discount", "promo_code", "final_price", "chef_id", "courier_id").
		From("orders").OrderBy("created_at DESC")

	if filter.StartAt != nil {
		builder = builder.Where(squirrel.GtOrEq{"created_at": *filter.StartAt})
	}
	if filter.EndAt != nil {
		builder = builder.Where(squirrel.LtOrEq{"created_at": *filter.EndAt})
	}
	if filter.Status != nil {
		builder = builder.Where(squirrel.Eq{"status": *filter.Status})
	}
	if filter.CustomerID != nil {
		builder = builder.Where(squirrel.Eq{"customer_id": *filter.CustomerID})
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build filtered query: %w", err)
	}

	rows, err := db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query filtered orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	var orderIDs []string

	for rows.Next() {
		var (
			oid, number, cid, promo, city, street, house, apartment, floor, entrance, comment string
			status                                                                            int
			createdAt                                                                         time.Time
			delPrice, discount, finalPrice                                                    common.Money
			chefID, courierID                                                                 sql.NullString
		)
		if err := rows.Scan(&oid, &number, &cid, &status, &createdAt, &city, &street, &house, &apartment, &floor, &entrance, &comment, &delPrice, &discount, &promo, &finalPrice, &chefID, &courierID); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		addr := domain.DeliveryAddress{City: city, Street: street, House: house, Apartment: apartment, Floor: floor, Entrance: entrance, Comment: comment}
		o := domain.ReconstructOrder(oid, number, cid, domain.OrderStatus(status), createdAt, addr, delPrice, discount, promo, finalPrice, nil, chefID.String, courierID.String)
		orders = append(orders, o)
		orderIDs = append(orderIDs, oid)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(orderIDs) == 0 {
		return orders, nil
	}

	itemsSql, itemsArgs, err := r.sb.Select("i.id", "i.order_id", "i.product_id", "i.product_name", "i.quantity", "i.base_price", "i.size_multiplier", "t.name", "t.price").
		From("order_items i").
		LeftJoin("order_item_toppings t ON i.id = t.order_item_id").
		Where(squirrel.Eq{"i.order_id": orderIDs}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build items query: %w", err)
	}

	itemRows, err := db.Query(ctx, itemsSql, itemsArgs...)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}
	defer itemRows.Close()

	type itemData struct {
		item    *domain.OrderItem
		orderID string
	}
	itemsByOrder := make(map[string][]*domain.OrderItem)
	itemsMap := make(map[string]*itemData)
	var itemsList []*itemData

	for itemRows.Next() {
		var (
			iid, oid, pid, pname string
			qty                  int
			base                 common.Money
			size                 float64
			tname                sql.NullString
			tprice               *common.Money
		)
		if err := itemRows.Scan(&iid, &oid, &pid, &pname, &qty, &base, &size, &tname, &tprice); err != nil {
			return nil, fmt.Errorf("scan item/topping: %w", err)
		}

		iData, ok := itemsMap[iid]
		if !ok {
			item := domain.ReconstructOrderItem(iid, pid, pname, qty, base, size, nil)
			iData = &itemData{item: item, orderID: oid}
			itemsMap[iid] = iData
			itemsList = append(itemsList, iData)
		}

		if tname.Valid && tprice != nil {
			iData.item.AddReconstructedTopping(tname.String, *tprice)
		}
	}

	// Preserve item order
	for _, iData := range itemsList {
		itemsByOrder[iData.orderID] = append(itemsByOrder[iData.orderID], iData.item)
	}

	var result []*domain.Order
	for _, o := range orders {
		items := itemsByOrder[o.ID()]
		if items == nil {
			items = make([]*domain.OrderItem, 0)
		}
		addr := o.Address()
		fullOrder := domain.ReconstructOrder(o.ID(), o.OrderNumber(), o.CustomerID(), o.Status(), o.CreatedAt(), addr, o.DeliveryPrice(), o.Discount(), o.PromoCode(), o.FinalPrice(), items, o.ChefID(), o.CourierID())
		result = append(result, fullOrder)
	}

	return result, nil
}


// --- Outbox ---

func (r *orderRepo) SaveOutboxEvent(ctx context.Context, event *domain.OutboxEvent) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Insert("outbox_events").
		Columns("id", "type", "payload", "created_at").
		Values(event.ID, event.Type, event.Payload, event.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("build outbox query: %w", err)
	}

	if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("exec save outbox event: %w", err)
	}
	return nil
}

// --- Promo Codes ---

func (r *orderRepo) SavePromo(ctx context.Context, p *domain.PromoCode) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	r.log.Info("Saving promo code", slog.String("code", p.Code()))
	sqlStr, args, err := r.sb.Insert("promo_codes").
		Columns("id", "code", "discount_type", "discount_amount", "is_active").
		Values(p.ID(), p.Code(), p.DiscountType(), p.DiscountAmount(), p.IsActive()).
		Suffix("ON CONFLICT (code) DO UPDATE SET discount_amount = EXCLUDED.discount_amount, is_active = EXCLUDED.is_active").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}
	if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("exec save promo: %w", err)
	}
	return nil
}

func (r *orderRepo) FindPromoByCode(ctx context.Context, code string) (*domain.PromoCode, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	var (
		id, c, dType string
		amount       common.Money
		active       bool
		expires      sql.NullTime
	)
	err := db.QueryRow(ctx, "SELECT id, code, discount_type, discount_amount, is_active, expires_at FROM promo_codes WHERE code = $1", code).
		Scan(&id, &c, &dType, &amount, &active, &expires)
	if err != nil {
		return nil, fmt.Errorf("query promo: %w", err)
	}
	return domain.NewPromoCode(id, c, dType, amount, active, expires.Time), nil
}

func (r *orderRepo) ListPromos(ctx context.Context) ([]*domain.PromoCode, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	rows, err := db.Query(ctx, "SELECT id, code, discount_type, discount_amount, is_active, expires_at FROM promo_codes ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("query promos: %w", err)
	}
	defer rows.Close()

	var res []*domain.PromoCode
	for rows.Next() {
		var (
			id, c, dType string
			amount       common.Money
			active       bool
			expires      sql.NullTime
		)
		if err := rows.Scan(&id, &c, &dType, &amount, &active, &expires); err != nil {
			r.log.Error("Failed to scan promo", slog.Any("error", err))
			continue
		}
		res = append(res, domain.NewPromoCode(id, c, dType, amount, active, expires.Time))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return res, nil
}

func (r *orderRepo) DeletePromo(ctx context.Context, id string) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	if _, err := db.Exec(ctx, "DELETE FROM promo_codes WHERE id = $1", id); err != nil {
		return fmt.Errorf("exec delete: %w", err)
	}
	return nil
}
