package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/catalog"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type CatalogUseCase struct {
	repo catalog.ProductRepository
}

func NewCatalogUseCase(repo catalog.ProductRepository) *CatalogUseCase {
	return &CatalogUseCase{repo: repo}
}

func (uc *CatalogUseCase) UpdatePrice(ctx context.Context, productID string, newPrice common.Money) error {
	if productID == "" {
		return fmt.Errorf("%w: product ID is required", ErrInvalidInput)
	}

	product, err := uc.repo.FindByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to find product %s: %w", productID, err)
	}

	if err := product.UpdatePrice(newPrice); err != nil {
		return fmt.Errorf("invalid price update: %w", err)
	}

	if err := uc.repo.Save(ctx, product); err != nil {
		return fmt.Errorf("failed to persist price update: %w", err)
	}

	return nil
}

func (uc *CatalogUseCase) SetAvailability(ctx context.Context, productID string, available bool) error {
	if productID == "" {
		return fmt.Errorf("%w: product ID is required", ErrInvalidInput)
	}

	product, err := uc.repo.FindByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to find product %s: %w", productID, err)
	}

	product.SetAvailability(available)

	if err := uc.repo.Save(ctx, product); err != nil {
		return fmt.Errorf("failed to persist availability update: %w", err)
	}

	return nil
}

func (uc *CatalogUseCase) CreateProduct(ctx context.Context, name, desc string, cat catalog.CategoryType, price common.Money) (*catalog.Product, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: product name is required", ErrInvalidInput)
	}

	product, err := catalog.NewProduct(name, desc, cat, price)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize product: %w", err)
	}

	if err := uc.repo.Save(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to save new product: %w", err)
	}

	return product, nil
}
