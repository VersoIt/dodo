package catalog

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// --- Value Objects ---

// Money - денежная единица (чтобы не путать с просто числами).
// В реальном проде лучше int64 (копейки), но для соответствия формулам ТЗ оставим float64.
type Money float64

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
// Может существовать отдельно от продукта.
type Ingredient struct {
	ID      string
	Name    string
	Cost    Money  // cost_i
	Unit    string // гр, мл
	InStock bool
}

// IngredientRef - Связь продукта с ингредиентом (Value Object внутри Product).
type IngredientRef struct {
	IngredientID string
	Quantity     float64
	IsRemovable  bool
}

// --- Aggregate ---

// Product - Агрегат товара.
// Поля скрыты (unexported), доступ только через методы. Это защищает инварианты.
type Product struct {
	id          string
	name        string
	description string
	category    CategoryType
	basePrice   Money           // base_i
	ingredients []IngredientRef // Рецептура
	imageUrl    string
	isAvailable bool
	createdAt   time.Time
}

// --- Factory ---

// NewProduct создает новый продукт, проверяя правила (Инварианты).
func NewProduct(name, desc string, cat CategoryType, basePrice Money) (*Product, error) {
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

// ID getter
func (p *Product) ID() string { return p.id }

// AddIngredient добавляет ингредиент в рецептуру.
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

// CalculatePrice рассчитывает цену с учетом модификатора размера (size_m).
// Формула (1) из ТЗ: base_i * size_m
// Топпинги считаются уже в заказе, так как они опциональны для конкретного инстанса еды.
func (p *Product) CalculateBasePrice(sizeMultiplier float64) Money {
	if sizeMultiplier <= 0 {
		sizeMultiplier = 1.0
	}
	return p.basePrice * Money(sizeMultiplier)
}

// MarkAvailable меняет доступность товара.
func (p *Product) MarkAvailable(available bool) {
	p.isAvailable = available
}

// IsAvailable getter
func (p *Product) IsAvailable() bool { return p.isAvailable }

// BasePrice getter
func (p *Product) BasePrice() Money { return p.basePrice }

// --- Repository Interface ---

type ProductRepository interface {
	FindAll() ([]*Product, error)
	FindByID(id string) (*Product, error)
	Save(p *Product) error
}

