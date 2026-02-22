package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/versoit/diploma/be/chat/internal/domain"
)

type messageRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	sb     squirrel.StatementBuilderType
	log    *slog.Logger
}

func NewMessageRepository(pool *pgxpool.Pool, log *slog.Logger) domain.MessageRepository {
	return &messageRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		sb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:    log,
	}
}

func (r *messageRepo) Save(ctx context.Context, msg *domain.Message) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	
	sqlStr, args, err := r.sb.Insert("chat_messages").
		Columns("order_id", "sender_id", "sender_name", "role", "text", "is_read").
		Values(msg.OrderID, msg.SenderID, msg.SenderName, msg.Role, msg.Text, msg.IsRead).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	err = db.QueryRow(ctx, sqlStr, args...).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec save message: %w", err)
	}
	return nil
}

func (r *messageRepo) GetHistory(ctx context.Context, orderID uuid.UUID, limit int) ([]domain.Message, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	
	sqlStr, args, err := r.sb.Select("id", "order_id", "sender_id", "sender_name", "role", "text", "created_at", "is_read").
		From("chat_messages").
		Where(squirrel.Eq{"order_id": orderID}).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query history: %w", err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(&msg.ID, &msg.OrderID, &msg.SenderID, &msg.SenderName, &msg.Role, &msg.Text, &msg.CreatedAt, &msg.IsRead); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, msg)
	}
	
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *messageRepo) GetAfterID(ctx context.Context, orderID uuid.UUID, afterID int64) ([]domain.Message, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	sqlStr, args, err := r.sb.Select("id", "order_id", "sender_id", "sender_name", "role", "text", "created_at", "is_read").
		From("chat_messages").
		Where(squirrel.Eq{"order_id": orderID}).
		Where(squirrel.Gt{"id": afterID}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query after id: %w", err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(&msg.ID, &msg.OrderID, &msg.SenderID, &msg.SenderName, &msg.Role, &msg.Text, &msg.CreatedAt, &msg.IsRead); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *messageRepo) MarkAsRead(ctx context.Context, messageID int64) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	sqlStr, args, err := r.sb.Update("chat_messages").
		Set("is_read", true).
		Where(squirrel.Eq{"id": messageID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	if _, err := db.Exec(ctx, sqlStr, args...); err != nil {
		return fmt.Errorf("exec mark as read: %w", err)
	}

	return nil
}
