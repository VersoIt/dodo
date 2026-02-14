package repository

import (
	"context"
	"time"
	"log/slog"

	"github.com/versoit/diploma/services/notification"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type notificationRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
	log  *slog.Logger
}

func NewNotificationRepository(pool *pgxpool.Pool, log *slog.Logger) notification.NotificationRepository {
	return &notificationRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:  log,
	}
}

func (r *notificationRepo) Save(ctx context.Context, n *notification.Notification) error {
	sqlStr, args, err := r.sb.Insert("notifications").
		Columns("id", "user_id", "channel", "title", "message", "status", "sent_at", "error_msg").
		Values(n.ID(), n.UserID(), string(n.Channel()), n.Title(), n.Message(), n.Status(), n.SentAt(), n.Error()).
		ToSql()
	if err != nil {
		return err
	}
	
	r.log.Debug("saving notification", slog.String("id", n.ID()), slog.String("user_id", n.UserID()))
	
	_, err = r.pool.Exec(ctx, sqlStr, args...)
	if err != nil {
		r.log.Error("failed to save notification", slog.Any("error", err), slog.String("id", n.ID()))
	}
	return err
}

func (r *notificationRepo) FindByUserID(ctx context.Context, userID string) ([]*notification.Notification, error) {
	sqlStr, args, err := r.sb.Select("id", "user_id", "channel", "title", "message", "status", "sent_at", "error_msg").
		From("notifications").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("sent_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sqlStr, args...)
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
