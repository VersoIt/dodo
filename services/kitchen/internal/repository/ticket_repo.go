package repository

import (
	"context"
	"fmt"
	"time"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/versoit/diploma/services/kitchen"
)

type ticketRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
	log  *slog.Logger
}

func NewTicketRepository(pool *pgxpool.Pool, log *slog.Logger) kitchen.TicketRepository {
	return &ticketRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:  log,
	}
}

func (r *ticketRepo) Save(ctx context.Context, t *kitchen.KitchenTicket) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	sqlStr, args, err := r.sb.Insert("kitchen_tickets").
		Columns("id", "order_id", "status", "created_at", "start_cooking_time", "ready_time").
		Values(t.ID(), t.OrderID(), t.Status(), t.CreatedAt(), t.StartTime(), t.ReadyTime()).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, start_cooking_time = EXCLUDED.start_cooking_time, ready_time = EXCLUDED.ready_time").
		ToSql()
	if err != nil {
		return err
	}

	r.log.Debug("saving kitchen ticket", slog.String("ticket_id", t.ID()), slog.Int("status", int(t.Status())))

	_, err = tx.Exec(ctx, sqlStr, args...)
	if err != nil {
		r.log.Error("failed to save kitchen ticket", slog.Any("error", err), slog.String("ticket_id", t.ID()))
		return err
	}

	_, _ = tx.Exec(ctx, "DELETE FROM kitchen_items WHERE ticket_id = $1", t.ID())
	for _, item := range t.Items() {
		sqlStr, args, _ = r.sb.Insert("kitchen_items").
			Columns("ticket_id", "product_id", "name", "quantity", "comment").
			Values(t.ID(), item.ProductID, item.Name, item.Quantity, item.Comment).
			Suffix("RETURNING id").
			ToSql()

		var itemID int64
		err = tx.QueryRow(ctx, sqlStr, args...).Scan(&itemID)
		if err != nil {
			return err
		}

		for _, ing := range item.Ingredients {
			sqlStr, args, _ = r.sb.Insert("kitchen_item_ingredients").
				Columns("kitchen_item_id", "ingredient_name").
				Values(itemID, ing).
				ToSql()
			_, err = tx.Exec(ctx, sqlStr, args...)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *ticketRepo) FindByID(ctx context.Context, id string) (*kitchen.KitchenTicket, error) {
	sqlStr, args, err := r.sb.Select("id", "order_id", "status", "created_at", "start_cooking_time", "ready_time").
		From("kitchen_tickets").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	return r.scanTicket(ctx, r.pool.QueryRow(ctx, sqlStr, args...))
}

func (r *ticketRepo) FindPending(ctx context.Context) ([]*kitchen.KitchenTicket, error) {
	sqlStr, args, err := r.sb.Select("id", "order_id", "status", "created_at", "start_cooking_time", "ready_time").
		From("kitchen_tickets").
		Where(squirrel.NotEq{"status": kitchen.TicketReady}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*kitchen.KitchenTicket
	for rows.Next() {
		t, err := r.scanTicket(ctx, rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (r *ticketRepo) scanTicket(ctx context.Context, row pgx.Row) (*kitchen.KitchenTicket, error) {
	var (
		id, oid string
		status  int
		cat     time.Time
		st, rdy *time.Time
	)

	if err := row.Scan(&id, &oid, &status, &cat, &st, &rdy); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("ticket not found")
		}
		return nil, err
	}

	itemSql, itemArgs, _ := r.sb.Select("id", "product_id", "name", "quantity", "comment").
		From("kitchen_items").
		Where(squirrel.Eq{"ticket_id": id}).
		ToSql()

	irows, err := r.pool.Query(ctx, itemSql, itemArgs...)
	if err != nil {
		return nil, err
	}
	defer irows.Close()

	var items []kitchen.KitchenItem
	for irows.Next() {
		var (
			itemID           int64
			pid, name, comm string
			qty              int
		)
		if err := irows.Scan(&itemID, &pid, &name, &qty, &comm); err != nil {
			return nil, err
		}

		ingSql, ingArgs, _ := r.sb.Select("ingredient_name").
			From("kitchen_item_ingredients").
			Where(squirrel.Eq{"kitchen_item_id": itemID}).
			ToSql()

		ingRows, err := r.pool.Query(ctx, ingSql, ingArgs...)
		if err != nil {
			return nil, err
		}
		
		var ingredients []string
		for ingRows.Next() {
			var ingName string
			if err := ingRows.Scan(&ingName); err != nil {
				ingRows.Close()
				return nil, err
			}
			ingredients = append(ingredients, ingName)
		}
		ingRows.Close()

		items = append(items, kitchen.KitchenItem{
			ProductID:   pid,
			Name:        name,
			Ingredients: ingredients,
			Quantity:    qty,
			Comment:     comm,
		})
	}

	var stVal, rdyVal time.Time
	if st != nil {
		stVal = *st
	}
	if rdy != nil {
		rdyVal = *rdy
	}

	return kitchen.ReconstructTicket(id, oid, kitchen.TicketStatus(status), cat, stVal, rdyVal, items), nil
}
