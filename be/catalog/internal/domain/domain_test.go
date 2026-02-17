package domain

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestNewProduct(t *testing.T) {
	name, _ := NewProductName("Pizza")
	price, _ := NewPrice(500)
	p, err := NewProduct(name, "Tasty", CatClassic, price, "http://img.jpg")
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	if p.Name() != "Pizza" {
		t.Errorf("expected Pizza, got %s", p.Name())
	}
	if !p.BasePrice().Equal(decimal.NewFromInt(500)) {
		t.Errorf("expected 500, got %v", p.BasePrice())
	}
}

func TestProduct_UpdatePrice(t *testing.T) {
	name, _ := NewProductName("Pizza")
	price, _ := NewPrice(500)
	p, _ := NewProduct(name, "", CatClassic, price, "")

	newPrice, _ := NewPrice(600)
	p.UpdatePrice(newPrice)
	if !p.BasePrice().Equal(decimal.NewFromInt(600)) {
		t.Errorf("failed to update price")
	}

	_, err := NewPrice(-1)
	if err != ErrNegativePrice {
		t.Errorf("expected ErrNegativePrice, got %v", err)
	}
}

func TestProduct_AddIngredient(t *testing.T) {
	name, _ := NewProductName("Pizza")
	price, _ := NewPrice(500)
	p, _ := NewProduct(name, "", CatClassic, price, "")
	err := p.AddIngredient("ing-1", 10.5, true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ingredients := p.Ingredients()
	if len(ingredients) != 1 || ingredients[0].IngredientID != "ing-1" {
		t.Errorf("ingredient not added correctly")
	}
}
