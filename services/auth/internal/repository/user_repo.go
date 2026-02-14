package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/versoit/diploma/services/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type userRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
	log  *slog.Logger
}

func NewUserRepository(pool *pgxpool.Pool, log *slog.Logger) auth.UserRepository {
	return &userRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:  log,
	}
}

func (r *userRepo) Save(ctx context.Context, u *auth.User) error {
	sqlStr, args, err := r.sb.Insert("users").
		Columns("id", "email", "password_hash", "role", "is_client", "name", "phone", "bonus_points", "updated_at").
		Values(u.ID(), u.Email(), u.HashedPassword(), u.Role(), u.IsClient(), u.Name(), u.Phone(), u.BonusPoints(), squirrel.Expr("NOW()")).
		Suffix("ON CONFLICT (id) DO UPDATE SET email = EXCLUDED.email, password_hash = EXCLUDED.password_hash, role = EXCLUDED.role, name = EXCLUDED.name, phone = EXCLUDED.phone, bonus_points = EXCLUDED.bonus_points, updated_at = NOW()").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	r.log.Debug("saving user to db", slog.String("user_id", u.ID()), slog.String("email", u.Email()))

	_, err = r.pool.Exec(ctx, sqlStr, args...)
	if err != nil {
		r.log.Error("db error saving user", slog.Any("error", err), slog.String("user_id", u.ID()))
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	sqlStr, args, err := r.sb.Select("id", "email", "password_hash", "role", "is_client", "name", "phone", "bonus_points").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	return r.scanUser(r.pool.QueryRow(ctx, sqlStr, args...))
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*auth.User, error) {
	sqlStr, args, err := r.sb.Select("id", "email", "password_hash", "role", "is_client", "name", "phone", "bonus_points").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	return r.scanUser(r.pool.QueryRow(ctx, sqlStr, args...))
}

func (r *userRepo) scanUser(row pgx.Row) (*auth.User, error) {
	var (
		id, email, pwdHash string
		name, phone        sql.NullString
		role               int
		isClient           bool
		bonus              int
	)

	err := row.Scan(&id, &email, &pwdHash, &role, &isClient, &name, &phone, &bonus)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, auth.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return auth.ReconstructUser(id, email, pwdHash, auth.Role(role), isClient, name.String, phone.String, bonus), nil
}
