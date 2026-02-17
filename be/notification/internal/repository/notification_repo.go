package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/versoit/diploma/be/notification/internal/domain"
)

type notificationRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	sb     squirrel.StatementBuilderType
	log    *slog.Logger
}

func NewNotificationRepository(pool *pgxpool.Pool, log *slog.Logger) domain.NotificationRepository {
	return &notificationRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		sb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:    log,
	}
}

func (r *notificationRepo) Save(ctx context.Context, n *domain.Notification) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Insert("notifications").
		Columns("id", "user_id", "channel", "title", "message", "status", "sent_at", "error_msg").
		Values(n.ID(), n.UserID(), string(n.Channel()), n.Title(), n.Message(), n.Status(), n.SentAt(), n.Error()).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	r.log.Debug("saving notification", slog.String("id", n.ID()), slog.String("user_id", n.UserID()))

	if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("exec save notification: %w", err)
	}
	return nil
}

func (r *notificationRepo) FindByUserID(ctx context.Context, userID string) ([]*domain.Notification, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("id", "user_id", "channel", "title", "message", "status", "sent_at", "error_msg").
		From("notifications").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("sent_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query notifications: %w", err)
	}
	defer rows.Close()

	var notes []*domain.Notification
	for rows.Next() {
		var (
			id, uid, ch, title, msg, errStr, statusStr string
			at                                         time.Time
		)
		if err := rows.Scan(&id, &uid, &ch, &title, &msg, &statusStr, &at, &errStr); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		notes = append(notes, domain.ReconstructNotification(id, uid, domain.Channel(ch), title, msg, statusStr == "sent", at, errStr))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return notes, nil
}
