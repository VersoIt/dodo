package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/services/catalog/internal/domain"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type CatalogUseCase struct {
	repo domain.ProductRepository
	tm   trm.Manager
	log  *slog.Logger
}

func NewCatalogUseCase(repo domain.ProductRepository, tm trm.Manager, log *slog.Logger) *CatalogUseCase {
	return &CatalogUseCase{
		repo: repo,
		tm:   tm,
		log:  log,
	}
}

func (uc *CatalogUseCase) UpdatePrice(ctx context.Context, productID string, newPrice float64) error {
	if productID == "" {
		return fmt.Errorf("%w: product ID is required", ErrInvalidInput)
	}
	priceVO, err := domain.NewPrice(newPrice)
	if err != nil {
		return fmt.Errorf("%w: invalid price: %v", ErrInvalidInput, err)
	}

	return uc.tm.Do(ctx, func(ctx context.Context) error {
		product, err := uc.repo.FindByID(ctx, productID)
		if err != nil {
			return fmt.Errorf("failed to find product %s: %w", productID, err)
		}

		product.UpdatePrice(priceVO)

		if err = uc.repo.Save(ctx, product); err != nil {
			return fmt.Errorf("failed to persist price update: %w", err)
		}

		uc.log.Info("product price updated", slog.String("product_id", productID), slog.Float64("new_price", newPrice))
		return nil
	})
}

func (uc *CatalogUseCase) SetAvailability(ctx context.Context, productID string, available bool) error {
	if productID == "" {
		return fmt.Errorf("%w: product ID is required", ErrInvalidInput)
	}

	return uc.tm.Do(ctx, func(ctx context.Context) error {
		product, err := uc.repo.FindByID(ctx, productID)
		if err != nil {
			return fmt.Errorf("failed to find product %s: %w", productID, err)
		}

		product.SetAvailability(available)

		if err = uc.repo.Save(ctx, product); err != nil {
			return fmt.Errorf("failed to persist availability update: %w", err)
		}

		uc.log.Info("product availability updated", slog.String("product_id", productID), slog.Bool("available", available))
		return nil
	})
}

func (uc *CatalogUseCase) CreateProduct(ctx context.Context, name, desc string, cat domain.CategoryType, price float64, imageUrl string) (*domain.Product, error) {
	nameVO, err := domain.NewProductName(name)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid name: %v", ErrInvalidInput, err)
	}
	priceVO, err := domain.NewPrice(price)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid price: %v", ErrInvalidInput, err)
	}

	var product *domain.Product
	err = uc.tm.Do(ctx, func(ctx context.Context) error {
		var err error

		product, err = domain.NewProduct(nameVO, desc, cat, priceVO, imageUrl)
		if err != nil {
			return fmt.Errorf("failed to initialize product: %w", err)
		}

		if err = uc.repo.Save(ctx, product); err != nil {
			return fmt.Errorf("failed to save new product: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	uc.log.Info("new product created", slog.String("product_id", product.ID()), slog.String("name", name))
	return product, nil
}

func (uc *CatalogUseCase) UpdateProduct(ctx context.Context, id, name, desc string, cat domain.CategoryType, price float64, imageUrl string, isAvailable bool) (*domain.Product, error) {
	nameVO, err := domain.NewProductName(name)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid name: %v", ErrInvalidInput, err)
	}
	priceVO, err := domain.NewPrice(price)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid price: %v", ErrInvalidInput, err)
	}

	var product *domain.Product
	err = uc.tm.Do(ctx, func(ctx context.Context) error {
		var err error

		product, err = uc.repo.FindByID(ctx, id)
		if err != nil {
			return fmt.Errorf("find product: %w", err)
		}

		product.Update(nameVO, desc, cat, priceVO, imageUrl, isAvailable)

		if err = uc.repo.Save(ctx, product); err != nil {
			return fmt.Errorf("save product: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *CatalogUseCase) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	products, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch products: %w", err)
	}
	return products, nil
}

func (uc *CatalogUseCase) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	product, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, domain.ErrProductNotFound
		}

		return nil, fmt.Errorf("find product %s: %w", id, err)
	}

	return product, nil
}
