package outbox

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type Relay struct {
	pool *pgxpool.Pool
	nc   *nats.Conn
	log  *slog.Logger
}

func NewRelay(pool *pgxpool.Pool, nc *nats.Conn, log *slog.Logger) *Relay {
	return &Relay{pool: pool, nc: nc, log: log}
}

func (r *Relay) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.processBatch(ctx)
		}
	}
}

func (r *Relay) processBatch(ctx context.Context) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.log.Error("failed to begin outbox tx", "error", err)
		return
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT id, type, payload
		FROM outbox_events
		WHERE processed_at IS NULL
		ORDER BY created_at ASC
		LIMIT 50
		FOR UPDATE SKIP LOCKED
	`)
	if err != nil {
		r.log.Error("failed to query outbox", "error", err)
		return
	}
	defer rows.Close()

	var toUpdate []string

	for rows.Next() {
		var id, eventType string
		var payload []byte
		if err := rows.Scan(&id, &eventType, &payload); err != nil {
			r.log.Error("failed to scan outbox event", "error", err)
			continue
		}

		if err := r.nc.Publish(eventType, payload); err != nil {
			r.log.Error("failed to publish outbox event", "error", err, "id", id)
			continue
		}
		
		toUpdate = append(toUpdate, id)
	}
	rows.Close()

	if len(toUpdate) > 0 {
		_, err = tx.Exec(ctx, `
			UPDATE outbox_events
			SET processed_at = NOW()
			WHERE id = ANY($1)
		`, toUpdate)
		if err != nil {
			r.log.Error("failed to update outbox events", "error", err)
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error("failed to commit outbox tx", "error", err)
	}
}
