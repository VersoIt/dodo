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
	"github.com/versoit/diploma/services/kitchen/internal/domain"
)

type ticketRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	sb     squirrel.StatementBuilderType
	log    *slog.Logger
}

func NewTicketRepository(pool *pgxpool.Pool, log *slog.Logger) domain.TicketRepository {
	return &ticketRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		sb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:    log,
	}
}

func (r *ticketRepo) Save(ctx context.Context, t *domain.KitchenTicket) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Insert("kitchen_tickets").
		Columns("id", "order_id", "status", "created_at", "start_cooking_time", "ready_time").
		Values(t.ID(), t.OrderID(), t.Status(), t.CreatedAt(), t.StartTime(), t.ReadyTime()).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, start_cooking_time = EXCLUDED.start_cooking_time, ready_time = EXCLUDED.ready_time").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	r.log.Debug("saving kitchen ticket", slog.String("ticket_id", t.ID()), slog.Int("status", int(t.Status())))

	_, err = db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("exec save ticket: %w", err)
	}

	if _, err = db.Exec(ctx, "DELETE FROM kitchen_items WHERE ticket_id = $1", t.ID()); err != nil {
		return fmt.Errorf("delete old items: %w", err)
	}
	for _, item := range t.Items() {
		sqlStr, args, err = r.sb.Insert("kitchen_items").
			Columns("ticket_id", "product_id", "name", "quantity", "comment").
			Values(t.ID(), item.ProductID, item.Name, item.Quantity, item.Comment).
			Suffix("RETURNING id").
			ToSql()
		if err != nil {
			return fmt.Errorf("build item query: %w", err)
		}

		var itemID int64
		err = db.QueryRow(ctx, sqlStr, args...).Scan(&itemID)
		if err != nil {
			return fmt.Errorf("scan item id: %w", err)
		}

		for _, ing := range item.Ingredients {
			sqlStr, args, err = r.sb.Insert("kitchen_item_ingredients").
				Columns("kitchen_item_id", "ingredient_name").
				Values(itemID, ing).
				ToSql()
			if err != nil {
				return fmt.Errorf("build ingredient query: %w", err)
			}
			if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
				return fmt.Errorf("exec save ingredient: %w", err)
			}
		}
	}

	return nil
}

func (r *ticketRepo) FindByID(ctx context.Context, id string) (*domain.KitchenTicket, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("id", "order_id", "status", "created_at", "start_cooking_time", "ready_time").
		From("kitchen_tickets").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	return r.scanTicket(ctx, db.QueryRow(ctx, sqlStr, args...))
}

func (r *ticketRepo) FindPending(ctx context.Context) ([]*domain.KitchenTicket, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("id", "order_id", "status", "created_at", "start_cooking_time", "ready_time").
		From("kitchen_tickets").
		Where(squirrel.NotEq{"status": domain.TicketReady}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query pending: %w", err)
	}
	defer rows.Close()

	var tickets []*domain.KitchenTicket
	for rows.Next() {
		t, err := r.scanTicket(ctx, rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return tickets, nil
}

func (r *ticketRepo) scanTicket(ctx context.Context, row pgx.Row) (*domain.KitchenTicket, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	var (
		id, oid string
		status  int
		cat     time.Time
		st, rdy *time.Time
	)

	if err := row.Scan(&id, &oid, &status, &cat, &st, &rdy); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ticket not found")
		}
		return nil, fmt.Errorf("scan ticket: %w", err)
	}

	itemSql, itemArgs, err := r.sb.Select("id", "product_id", "name", "quantity", "comment").
		From("kitchen_items").
		Where(squirrel.Eq{"ticket_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build items query: %w", err)
	}

	// We need Query here, so we cast db.
	type query interface {
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	}
	q, ok := db.(query)
	if !ok {
		return nil, fmt.Errorf("db does not support Query")
	}

	irows, err := q.Query(ctx, itemSql, itemArgs...)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}
	defer irows.Close()

	var items []domain.KitchenItem
	for irows.Next() {
		var (
			itemID          int64
			pid, name, comm string
			qty             int
		)
		if err := irows.Scan(&itemID, &pid, &name, &qty, &comm); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}

		ingSql, ingArgs, err := r.sb.Select("ingredient_name").
			From("kitchen_item_ingredients").
			Where(squirrel.Eq{"kitchen_item_id": itemID}).
			ToSql()
		if err != nil {
			return nil, fmt.Errorf("build ingredients query: %w", err)
		}

		ingRows, err := q.Query(ctx, ingSql, ingArgs...)
		if err != nil {
			return nil, fmt.Errorf("query ingredients: %w", err)
		}

		var ingredients []string
		for ingRows.Next() {
			var ingName string
			if err := ingRows.Scan(&ingName); err != nil {
				ingRows.Close()
				return nil, fmt.Errorf("scan ingredient: %w", err)
			}
			ingredients = append(ingredients, ingName)
		}
		ingRows.Close()

		items = append(items, domain.KitchenItem{
			ProductID:   pid,
			Name:        name,
			Ingredients: ingredients,
			Quantity:    qty,
			Comment:     comm,
		})
	}
	if err := irows.Err(); err != nil {
		return nil, fmt.Errorf("items rows error: %w", err)
	}

	var stVal, rdyVal time.Time
	if st != nil {
		stVal = *st
	}
	if rdy != nil {
		rdyVal = *rdy
	}

	return domain.ReconstructTicket(id, oid, domain.TicketStatus(status), cat, stVal, rdyVal, items), nil
}
