package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/versoit/diploma/be/catalog/internal/domain"
)

type productRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	sb     squirrel.StatementBuilderType
	log    *slog.Logger
}

func NewProductRepository(pool *pgxpool.Pool, log *slog.Logger) domain.ProductRepository {
	return &productRepo{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		sb:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:    log,
	}
}

func (r *productRepo) FindAll(ctx context.Context) ([]*domain.Product, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("id", "name", "description", "category", "base_price", "image_url", "is_available").
		From("products").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query all: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var (
			pid, name string
			desc, img sql.NullString
			cat       int
			price     float64
			isAvail   bool
		)
		if err := rows.Scan(&pid, &name, &desc, &cat, &price, &img, &isAvail); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}

		ingredients, err := r.fetchIngredients(ctx, pid)
		if err != nil {
			return nil, fmt.Errorf("fetch ingredients: %w", err)
		}

		products = append(products, domain.ReconstructProduct(pid, name, desc.String, domain.CategoryType(cat), price, img.String, isAvail, ingredients))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return products, nil
}

func (r *productRepo) fetchIngredients(ctx context.Context, productID string) ([]domain.IngredientRef, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("ingredient_id", "quantity", "is_removable").
		From("product_ingredients").
		Where(squirrel.Eq{"product_id": productID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []domain.IngredientRef
	for rows.Next() {
		var (
			id  string
			qty float64
			rem bool
		)
		if err := rows.Scan(&id, &qty, &rem); err != nil {
			return nil, fmt.Errorf("scan ingredient: %w", err)
		}
		ingredients = append(ingredients, domain.IngredientRef{
			IngredientID: id,
			Quantity:     qty,
			IsRemovable:  rem,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return ingredients, nil
}

func (r *productRepo) Save(ctx context.Context, p *domain.Product) error {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Insert("products").
		Columns("id", "name", "description", "category", "base_price", "image_url", "is_available").
		Values(p.ID(), p.Name(), p.Description(), p.Category(), p.BasePrice().InexactFloat64(), p.ImageURL(), p.IsAvailable()).
		Suffix("ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, category = EXCLUDED.category, base_price = EXCLUDED.base_price, image_url = EXCLUDED.image_url, is_available = EXCLUDED.is_available").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	r.log.Debug("saving product", slog.String("product_id", p.ID()), slog.String("name", p.Name()))

	_, err = db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("exec save product: %w", err)
	}

	_, err = db.Exec(ctx, "DELETE FROM product_ingredients WHERE product_id = $1", p.ID())
	if err != nil {
		return fmt.Errorf("delete ingredients: %w", err)
	}

	for _, ing := range p.Ingredients() {
		sqlStr, args, err = r.sb.Insert("product_ingredients").
			Columns("product_id", "ingredient_id", "quantity", "is_removable").
			Values(p.ID(), ing.IngredientID, ing.Quantity, ing.IsRemovable).
			ToSql()
		if err != nil {
			return fmt.Errorf("build ing query: %w", err)
		}
		if _, err = db.Exec(ctx, sqlStr, args...); err != nil {
			return fmt.Errorf("exec save ingredient: %w", err)
		}
	}

	return nil
}

func (r *productRepo) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlStr, args, err := r.sb.Select("id", "name", "description", "category", "base_price", "image_url", "is_available").
		From("products").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var (
		pid, name string
		desc, img sql.NullString
		cat       int
		price     float64
		isAvail   bool
	)

	err = db.QueryRow(ctx, sqlStr, args...).Scan(&pid, &name, &desc, &cat, &price, &img, &isAvail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("query product: %w", err)
	}

	ingredients, err := r.fetchIngredients(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch ingredients: %w", err)
	}

	return domain.ReconstructProduct(pid, name, desc.String, domain.CategoryType(cat), price, img.String, isAvail, ingredients), nil
}
