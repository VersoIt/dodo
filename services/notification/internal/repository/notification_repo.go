package repository

import (
	"context"
	"time"

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
		Columns("id", "user_id", "channel", "title", "message", "status", "sent_at", "error_msg").
		Values(n.ID(), n.UserID(), string(n.Channel()), n.Title(), n.Message(), n.Status(), n.SentAt(), n.Error()).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *notificationRepo) FindByUserID(ctx context.Context, userID string) ([]*notification.Notification, error) {
	sql, args, err := r.sb.Select("id", "user_id", "channel", "title", "message", "status", "sent_at", "error_msg").
		From("notifications").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("sent_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*notification.Notification
	for rows.Next() {
		var (
			id, uid, ch, title, msg, errStr string
			status                          bool
			at                              time.Time
		)
		if err := rows.Scan(&id, &uid, &ch, &title, &msg, &status, &at, &errStr); err != nil {
			return nil, err
		}
		notes = append(notes, notification.ReconstructNotification(id, uid, notification.Channel(ch), title, msg, status, at, errStr))
	}
	return notes, nil
}