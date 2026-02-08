package repository

import (
	"context"

	"github.com/versoit/diploma/services/kitchen"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type ticketRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewTicketRepository(pool *pgxpool.Pool) kitchen.TicketRepository {
	return &ticketRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ticketRepo) Save(ctx context.Context, t *kitchen.KitchenTicket) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.sb.Insert("kitchen_tickets").
		Columns("id", "order_id", "status").
		Values(t.ID(), t.OrderID(), t.Status()).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status").
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *ticketRepo) FindByID(ctx context.Context, id string) (*kitchen.KitchenTicket, error) {
	// Simple implementation
	return nil, nil
}

func (r *ticketRepo) FindPending(ctx context.Context) ([]*kitchen.KitchenTicket, error) {
	return nil, nil
}