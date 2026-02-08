package repository

import (
	"context"

	"github.com/versoit/diploma/services/treasury"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type paymentRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewPaymentRepository(pool *pgxpool.Pool) treasury.PaymentRepository {
	return &paymentRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *paymentRepo) Save(ctx context.Context, p *treasury.Payment) error {
	sql, args, err := r.sb.Insert("payments").
		Columns("id", "order_id", "amount", "method", "status", "transaction_id").
		Values(p.ID(), p.OrderID(), p.Amount(), p.Method(), p.Status(), p.TransactionID()).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, transaction_id = EXCLUDED.transaction_id, updated_at = NOW()").
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *paymentRepo) FindByOrderID(ctx context.Context, orderID string) (*treasury.Payment, error) {
	// Implementation...
	return nil, nil
}