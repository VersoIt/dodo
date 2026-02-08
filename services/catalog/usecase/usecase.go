package usecase

import (
	"context"
	"fmt"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/catalog"
)

type CatalogUseCase struct {
	repo catalog.ProductRepository
}

func NewCatalogUseCase(repo catalog.ProductRepository) *CatalogUseCase {
	return &CatalogUseCase{repo: repo}
}

// UpdatePrice обновляет цену товара.
func (uc *CatalogUseCase) UpdatePrice(ctx context.Context, productID string, newPrice common.Money) error {
	product, err := uc.repo.FindByID(productID)
	if err != nil {
		return err
	}

	if err := product.UpdatePrice(newPrice); err != nil {
		return err
	}

	return uc.repo.Save(product)
}

// SetAvailability переключает статус доступности.
func (uc *CatalogUseCase) SetAvailability(ctx context.Context, productID string, available bool) error {
	product, err := uc.repo.FindByID(productID)
	if err != nil {
		return err
	}

	product.SetAvailability(available)

	return uc.repo.Save(product)
}

// CreateProduct создает новый товар в каталоге.
func (uc *CatalogUseCase) CreateProduct(ctx context.Context, name, desc string, cat catalog.CategoryType, price common.Money) (*catalog.Product, error) {
	product, err := catalog.NewProduct(name, desc, cat, price)
	if err != nil {
		return nil, fmt.Errorf("failed to create product entity: %w", err)
	}

	if err := uc.repo.Save(product); err != nil {
		return nil, fmt.Errorf("failed to save product: %w", err)
	}

	return product, nil
}
