package catalog

import (
	"testing"
)

func TestNewProduct(t *testing.T) {
	p, err := NewProduct("Pizza", "Tasty", CatClassic, 500)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	if p.Name() != "Pizza" {
		t.Errorf("expected Pizza, got %s", p.Name())
	}
	if p.BasePrice() != 500 {
		t.Errorf("expected 500, got %v", p.BasePrice())
	}
}

func TestProduct_UpdatePrice(t *testing.T) {
	p, _ := NewProduct("Pizza", "", CatClassic, 500)
	
	err := p.UpdatePrice(600)
	if err != nil || p.BasePrice() != 600 {
		t.Errorf("failed to update price")
	}

	if err := p.UpdatePrice(-1); err != ErrNegativePrice {
		t.Errorf("expected ErrNegativePrice, got %v", err)
	}
}

func TestProduct_AddIngredient(t *testing.T) {
	p, _ := NewProduct("Pizza", "", CatClassic, 500)
	err := p.AddIngredient("ing-1", 10.5, true)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ingredients := p.Ingredients()
	if len(ingredients) != 1 || ingredients[0].IngredientID != "ing-1" {
		t.Errorf("ingredient not added correctly")
	}
}
