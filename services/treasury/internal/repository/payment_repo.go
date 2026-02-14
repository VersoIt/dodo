package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"log/slog"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/treasury"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type paymentRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
	log  *slog.Logger
}

func NewPaymentRepository(pool *pgxpool.Pool, log *slog.Logger) treasury.PaymentRepository {
	return &paymentRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:  log,
	}
}

func (r *paymentRepo) Save(ctx context.Context, p *treasury.Payment) error {
	sqlStr, args, err := r.sb.Insert("payments").
		Columns("id", "order_id", "amount", "method", "status", "transaction_id").
		Values(p.ID(), p.OrderID(), p.Amount(), p.Method(), p.Status(), p.TransactionID()).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, transaction_id = EXCLUDED.transaction_id, updated_at = NOW()").
		ToSql()
	if err != nil {
		return err
	}

	r.log.Debug("saving payment", slog.String("payment_id", p.ID()), slog.String("order_id", p.OrderID()), slog.Int("status", int(p.Status())))

	_, err = r.pool.Exec(ctx, sqlStr, args...)
	if err != nil {
		r.log.Error("failed to save payment", slog.Any("error", err), slog.String("payment_id", p.ID()))
	}
	return err
}

func (r *paymentRepo) FindByOrderID(ctx context.Context, orderID string) (*treasury.Payment, error) {
	sqlStr, args, err := r.sb.Select("id", "order_id", "amount", "method", "status", "transaction_id", "created_at").
		From("payments").
		Where(squirrel.Eq{"order_id": orderID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		id, oid   string
		txid      sql.NullString
		amount    common.Money
		method    int
		status    int
		createdAt time.Time
	)

	err = r.pool.QueryRow(ctx, sqlStr, args...).Scan(&id, &oid, &amount, &method, &status, &txid, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, err
	}

	return treasury.ReconstructPayment(id, oid, amount, treasury.PaymentMethod(method), treasury.PaymentStatus(status), txid.String, createdAt), nil
}
