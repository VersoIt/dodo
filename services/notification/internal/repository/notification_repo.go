package repository

import (
	"context"

	"github.com/versoit/diploma/services/notification"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type notificationRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewNotificationRepository(pool *pgxpool.Pool) notification.NotificationRepository {
	return &notificationRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *notificationRepo) Save(ctx context.Context, n *notification.Notification) error {
	sql, args, err := r.sb.Insert("notifications").
		Columns("id", "user_id", "channel", "title", "message", "status").
		Values(n.ID(), n.UserID(), n.Channel(), n.Title(), n.Message(), n.Status()).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *notificationRepo) FindByUserID(ctx context.Context, userID string) ([]*notification.Notification, error) {
	return nil, nil
}