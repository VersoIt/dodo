package repository

import (
	"context"
	"log/slog"

	"github.com/versoit/diploma/services/catalog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type productRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
	log  *slog.Logger
}

func NewProductRepository(pool *pgxpool.Pool, log *slog.Logger) catalog.ProductRepository {
	return &productRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:  log,
	}
}

func (r *productRepo) FindAll(ctx context.Context) ([]*catalog.Product, error) {
	sqlStr, args, err := r.sb.Select("id", "name", "description", "category", "base_price", "image_url", "is_available").
		From("products").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*catalog.Product
	for rows.Next() {
		var (
			pid, name, desc, img string
			cat                  int
			price                float64
			isAvail              bool
		)
		if err := rows.Scan(&pid, &name, &desc, &cat, &price, &img, &isAvail); err != nil {
			return nil, err
		}

		ingredients, err := r.fetchIngredients(ctx, pid)
		if err != nil {
			return nil, err
		}

		products = append(products, catalog.ReconstructProduct(pid, name, desc, catalog.CategoryType(cat), price, img, isAvail, ingredients))
	}
	return products, nil
}

func (r *productRepo) fetchIngredients(ctx context.Context, productID string) ([]catalog.IngredientRef, error) {
	sqlStr, args, err := r.sb.Select("ingredient_id", "quantity", "is_removable").
		From("product_ingredients").
		Where(squirrel.Eq{"product_id": productID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []catalog.IngredientRef
	for rows.Next() {
		var (
			id   string
			qty  float64
			rem  bool
		)
		if err := rows.Scan(&id, &qty, &rem); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, catalog.IngredientRef{
			IngredientID: id,
			Quantity:     qty,
			IsRemovable:  rem,
		})
	}
	return ingredients, nil
}

func (r *productRepo) Save(ctx context.Context, p *catalog.Product) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	sqlStr, args, err := r.sb.Insert("products").
		Columns("id", "name", "description", "category", "base_price", "image_url", "is_available").
		Values(p.ID(), p.Name(), p.Description(), p.Category(), p.BasePrice().InexactFloat64(), p.ImageURL(), p.IsAvailable()).
		Suffix("ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, category = EXCLUDED.category, base_price = EXCLUDED.base_price, image_url = EXCLUDED.image_url, is_available = EXCLUDED.is_available").
		ToSql()
	if err != nil {
		return err
	}

	r.log.Debug("saving product", slog.String("product_id", p.ID()), slog.String("name", p.Name()))

	_, err = tx.Exec(ctx, sqlStr, args...)
	if err != nil {
		r.log.Error("failed to save product", slog.Any("error", err), slog.String("product_id", p.ID()))
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM product_ingredients WHERE product_id = $1", p.ID())
	if err != nil {
		return err
	}

	for _, ing := range p.Ingredients() {
		sqlStr, args, err := r.sb.Insert("product_ingredients").
			Columns("product_id", "ingredient_id", "quantity", "is_removable").
			Values(p.ID(), ing.IngredientID, ing.Quantity, ing.IsRemovable).
			ToSql()
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, sqlStr, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *productRepo) FindByID(ctx context.Context, id string) (*catalog.Product, error) {
	sqlStr, args, err := r.sb.Select("id", "name", "description", "category", "base_price", "image_url", "is_available").
		From("products").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		pid, name, desc, img string
		cat                  int
		price                float64 
		isAvail              bool
	)

	err = r.pool.QueryRow(ctx, sqlStr, args...).Scan(&pid, &name, &desc, &cat, &price, &img, &isAvail)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, catalog.ErrProductNotFound
		}
		return nil, err
	}

	ingredients, err := r.fetchIngredients(ctx, id)
	if err != nil {
		return nil, err
	}

	return catalog.ReconstructProduct(pid, name, desc, catalog.CategoryType(cat), price, img, isAvail, ingredients), nil
}
