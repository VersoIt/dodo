package catalog

import (
	"github.com/versoit/diploma/pkg/common"
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

// --- Errors ---

var (
	ErrInvalidDetails  = errors.New("invalid product details")
	ErrNegativePrice   = errors.New("price cannot be negative")
	ErrNegativeQty     = errors.New("quantity must be positive")
	ErrProductNotFound = errors.New("product not found")
)

// --- Factory ---

func NewProduct(name, desc string, cat CategoryType, basePrice common.Money) (*Product, error) {
	if name == "" {
		return nil, ErrInvalidDetails
	}
	if basePrice < 0 {
		return nil, ErrNegativePrice
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
		return ErrNegativeQty
	}
	p.ingredients = append(p.ingredients, IngredientRef{
		IngredientID: ingID,
		Quantity:     qty,
		IsRemovable:  removable,
	})
	return nil
}

func (p *Product) UpdatePrice(newPrice common.Money) error {
	if newPrice < 0 {
		return ErrNegativePrice
	}
	p.basePrice = newPrice
	return nil
}

func (p *Product) SetAvailability(available bool) {
	p.isAvailable = available
}

// --- Getters (Accessors) ---

func (p *Product) ID() string               { return p.id }
func (p *Product) Name() string             { return p.name }
func (p *Product) Description() string      { return p.description }
func (p *Product) Category() CategoryType   { return p.category }
func (p *Product) BasePrice() common.Money  { return p.basePrice }
func (p *Product) ImageURL() string         { return p.imageUrl }
func (p *Product) IsAvailable() bool        { return p.isAvailable }

// Ingredients возвращает КОПИЮ списка ингредиентов для защиты внутреннего состояния.
func (p *Product) Ingredients() []IngredientRef {
	result := make([]IngredientRef, len(p.ingredients))
	copy(result, p.ingredients)
	return result
}

// --- Repository Interface ---

type ProductRepository interface {
	FindAll() ([]*Product, error)
	FindByID(id string) (*Product, error)
	Save(p *Product) error
}

