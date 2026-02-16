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
	"github.com/versoit/diploma/services/treasury/internal/domain"
)

type paymentRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	sb     squirrel.StatementBuilderType
	log    *slog.Logger
}

func NewPaymentRepository(pool *pgxpool.Pool, log *slog.Logger) domain.PaymentRepository {
	return &paymentRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		sb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:    log,
	}
}

func (r *paymentRepo) Save(ctx context.Context, p *domain.Payment) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Insert("payments").
		Columns("id", "order_id", "amount", "method", "status", "transaction_id").
		Values(p.ID(), p.OrderID(), p.Amount(), p.Method(), p.Status(), p.TransactionID()).
		Suffix("ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, transaction_id = EXCLUDED.transaction_id, updated_at = NOW()").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	r.log.Debug("saving payment", slog.String("payment_id", p.ID()), slog.String("order_id", p.OrderID()), slog.Int("status", int(p.Status())))

	if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("exec save payment: %w", err)
	}
	return nil
}

func (r *paymentRepo) FindByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("id", "order_id", "amount", "method", "status", "transaction_id", "created_at").
		From("payments").
		Where(squirrel.Eq{"order_id": orderID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var (
		id, oid   string
		txid      sql.NullString
		amount    common.Money
		method    int
		status    int
		createdAt time.Time
	)

	err = db.QueryRow(ctx, sqlStr, args...).Scan(&id, &oid, &amount, &method, &status, &txid, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("query payment: %w", err)
	}

	return domain.ReconstructPayment(id, oid, amount, domain.PaymentMethod(method), domain.PaymentStatus(status), txid.String, createdAt), nil
}
