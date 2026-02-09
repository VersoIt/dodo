package repository

import (
	"context"

	"github.com/versoit/diploma/services/catalog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
)

type productRepo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewProductRepository(pool *pgxpool.Pool) catalog.ProductRepository {
	return &productRepo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *productRepo) FindAll(ctx context.Context) ([]*catalog.Product, error) {
	sql, args, err := r.sb.Select("id", "name", "description", "category", "base_price", "is_available").
		From("products").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*catalog.Product
	for rows.Next() {
		var (
			pid, name, desc string
			cat             int
			price           float64
			isAvail         bool
		)
		if err := rows.Scan(&pid, &name, &desc, &cat, &price, &isAvail); err != nil {
			return nil, err
		}

		ingredients, err := r.fetchIngredients(ctx, pid)
		if err != nil {
			return nil, err
		}

		products = append(products, catalog.ReconstructProduct(pid, name, desc, catalog.CategoryType(cat), price, isAvail, ingredients))
	}
	return products, nil
}

func (r *productRepo) fetchIngredients(ctx context.Context, productID string) ([]catalog.IngredientRef, error) {
	sql, args, err := r.sb.Select("ingredient_id", "quantity", "is_removable").
		From("product_ingredients").
		Where(squirrel.Eq{"product_id": productID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
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
	defer tx.Rollback(ctx)

	sql, args, err := r.sb.Insert("products").
		Columns("id", "name", "description", "category", "base_price", "is_available").
		Values(p.ID(), p.Name(), p.Description(), p.Category(), p.BasePrice(), p.IsAvailable()).
		Suffix("ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, category = EXCLUDED.category, base_price = EXCLUDED.base_price, is_available = EXCLUDED.is_available").
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	// Simple way: delete and insert ingredients
	_, err = tx.Exec(ctx, "DELETE FROM product_ingredients WHERE product_id = $1", p.ID())
	if err != nil {
		return err
	}

	for _, ing := range p.Ingredients() {
		sql, args, err := r.sb.Insert("product_ingredients").
			Columns("product_id", "ingredient_id", "quantity", "is_removable").
			Values(p.ID(), ing.IngredientID, ing.Quantity, ing.IsRemovable).
			ToSql()
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *productRepo) FindByID(ctx context.Context, id string) (*catalog.Product, error) {
	sql, args, err := r.sb.Select("id", "name", "description", "category", "base_price", "is_available").
		From("products").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		pid, name, desc string
		cat             int
		price           float64 // Need to handle Money
		isAvail         bool
	)

	// Note: Money scanning might need a custom scanner if not using float64 or string
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&pid, &name, &desc, &cat, &price, &isAvail)
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

	return catalog.ReconstructProduct(pid, name, desc, catalog.CategoryType(cat), price, isAvail, ingredients), nil
}