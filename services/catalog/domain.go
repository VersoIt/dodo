package catalog

import (
	"diploma/pkg/common"
	"errors"
	"time"

	"github.com/google/uuid"
)

// --- Value Objects ---

// CategoryType - коды категорий
type CategoryType int

const (
	CatClassic    CategoryType = 0
	CatPremium    CategoryType = 1
	CatVegetarian CategoryType = 2
	CatSpicy      CategoryType = 3
	CatDrinks     CategoryType = 4
	CatDesserts   CategoryType = 5
)

// --- Entities ---

// Ingredient - Сущность ингредиента.
type Ingredient struct {
	ID      string
	Name    string
	Cost    common.Money
	Unit    string
	InStock bool
}

// IngredientRef - Связь продукта с ингредиентом.
type IngredientRef struct {
	IngredientID string
	Quantity     float64
	IsRemovable  bool
}

// --- Aggregate ---

// Product - Агрегат товара.
type Product struct {
	id          string
	name        string
	description string
	category    CategoryType
	basePrice   common.Money
	ingredients []IngredientRef
	imageUrl    string
	isAvailable bool
	createdAt   time.Time
}

// --- Factory ---

func NewProduct(name, desc string, cat CategoryType, basePrice common.Money) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if basePrice < 0 {
		return nil, errors.New("price cannot be negative")
	}

	return &Product{
		id:          uuid.New().String(),
		name:        name,
		description: desc,
		category:    cat,
		basePrice:   basePrice,
		ingredients: make([]IngredientRef, 0),
		isAvailable: true,
		createdAt:   time.Now(),
	}, nil
}

// --- Behavior (Methods) ---

func (p *Product) AddIngredient(ingID string, qty float64, removable bool) error {
	if qty <= 0 {
		return errors.New("quantity must be positive")
	}
	p.ingredients = append(p.ingredients, IngredientRef{
		IngredientID: ingID,
		Quantity:     qty,
		IsRemovable:  removable,
	})
	return nil
}

func (p *Product) CalculateBasePrice(sizeMultiplier float64) common.Money {
	if sizeMultiplier <= 0 {
		sizeMultiplier = 1.0
	}
	return p.basePrice * common.Money(sizeMultiplier)
}

func (p *Product) MarkAvailable(available bool) {
	p.isAvailable = available
}

// --- Getters (Accessors) ---
// Необходимы для маппинга в DTO или сохранения в БД

func (p *Product) ID() string                { return p.id }
func (p *Product) Name() string              { return p.name }
func (p *Product) Description() string       { return p.description }
func (p *Product) Category() CategoryType    { return p.category }
func (p *Product) BasePrice() common.Money   { return p.basePrice }
func (p *Product) Ingredients() []IngredientRef { return p.ingredients } // Возвращаем копию слайса лучше, но пока ок
func (p *Product) IsAvailable() bool         { return p.isAvailable }
func (p *Product) CreatedAt() time.Time      { return p.createdAt }

// --- Repository Interface ---

type ProductRepository interface {
	FindAll() ([]*Product, error)
	FindByID(id string) (*Product, error)
	Save(p *Product) error
}

