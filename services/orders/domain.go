package orders

import (
	"errors"
	"fmt"
	"time"

	"diploma/pkg/common"

	"github.com/google/uuid"
)

// --- Value Objects ---

type OrderStatus int

const (
	StatusCreated    OrderStatus = 0
	StatusPaid       OrderStatus = 1
	StatusCooking    OrderStatus = 2
	StatusReady      OrderStatus = 3
	StatusDelivering OrderStatus = 4
	StatusCompleted  OrderStatus = 5
	StatusCanceled   OrderStatus = 6
)

type DeliveryAddress struct {
	City      string
	Street    string
	House     string
	Apartment string
	Floor     string
	Comment   string
}

type Topping struct {
	Name  string
	Price common.Money
}

func (s OrderStatus) String() string {
	switch s {
	case StatusCreated:
		return "created"
	case StatusPaid:
		return "paid"
	case StatusCooking:
		return "cooking"
	case StatusReady:
		return "ready"
	case StatusDelivering:
		return "delivering"
	case StatusCompleted:
		return "completed"
	case StatusCanceled:
		return "canceled"
	default:
		return "unknown"
	}
}

// --- Entities ---

// OrderItem - Позиция заказа.
type OrderItem struct {
	productID      string
	productName    string
	quantity       int
	basePrice      common.Money // Базовая цена товара (без размера)
	sizeMultiplier float64      // Модификатор размера
	toppings       []Topping
}

func (i *OrderItem) CalculateTotal() common.Money {
	// Цена с учетом размера
	sizedPrice := i.basePrice * common.Money(i.sizeMultiplier)

	var toppingsPrice common.Money
	for _, t := range i.toppings {
		toppingsPrice += t.Price
	}

	// (Base*Size + Toppings) * Qty
	unitPrice := sizedPrice + toppingsPrice
	return unitPrice * common.Money(i.quantity)
}

// Getters for Item
func (i *OrderItem) ProductID() string       { return i.productID }
func (i *OrderItem) ProductName() string     { return i.productName }
func (i *OrderItem) Quantity() int           { return i.quantity }
func (i *OrderItem) BasePrice() common.Money { return i.basePrice }
func (i *OrderItem) Size() float64           { return i.sizeMultiplier }
func (i *OrderItem) Toppings() []Topping     { return i.toppings }

// --- Aggregate Root ---

type Order struct {
	id          string
	orderNumber string
	customerID  string
	status      OrderStatus
	createdAt   time.Time

	items   []*OrderItem
	address DeliveryAddress

	deliveryPrice common.Money
	discount      common.Money
	promoCode     string

	finalPrice common.Money
}

// --- Factory ---

func NewOrder(customerID string, address DeliveryAddress) *Order {
	return &Order{
		id:          uuid.New().String(),
		orderNumber: generateOrderNumber(),
		customerID:  customerID,
		status:      StatusCreated,
		createdAt:   time.Now(),
		address:     address,
		items:       make([]*OrderItem, 0),
	}
}

// --- Errors ---

var (
	ErrOrderLocked       = errors.New("order is locked for changes")
	ErrInvalidQty        = errors.New("quantity must be positive")
	ErrInvalidDiscount   = errors.New("invalid discount")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrOrderNotFound     = errors.New("order not found")
)

// ...

// --- Business Logic ---

func (o *Order) AddItem(productID, name string, qty int, productBasePrice common.Money, sizeMult float64, toppings []Topping) error {
	if o.status != StatusCreated {
		return ErrOrderLocked
	}
	if qty <= 0 {
		return ErrInvalidQty
	}

	// Defensive copy of toppings
	toppingsCopy := make([]Topping, len(toppings))
	copy(toppingsCopy, toppings)

	o.items = append(o.items, &OrderItem{
		productID:      productID,
		productName:    name,
		quantity:       qty,
		basePrice:      productBasePrice,
		sizeMultiplier: sizeMult,
		toppings:       toppingsCopy,
	})

	o.recalculate()
	return nil
}

func (o *Order) ApplyPromoCode(code string, discountAmount common.Money) error {
	if o.status != StatusCreated {
		return ErrOrderLocked
	}
	if discountAmount < 0 {
		return ErrInvalidDiscount
	}

	o.promoCode = code
	o.discount = discountAmount
	o.recalculate()
	return nil
}

func (o *Order) SetDeliveryPrice(price common.Money) {
	o.deliveryPrice = price
	o.recalculate()
}

func (o *Order) recalculate() {
	var total common.Money
	for _, item := range o.items {
		total += item.CalculateTotal()
	}

	o.finalPrice = total + o.deliveryPrice - o.discount
	if o.finalPrice < 0 {
		o.finalPrice = 0
	}
}

// --- State Machine ---

func (o *Order) MarkPaid() error {
	if o.status != StatusCreated {
		return fmt.Errorf("%w: cannot pay for order in status %s", ErrInvalidTransition, o.status)
	}
	o.status = StatusPaid
	return nil
}

func (o *Order) SendToKitchen() error {
	if o.status != StatusPaid {
		return fmt.Errorf("%w: order must be paid", ErrInvalidTransition)
	}
	o.status = StatusCooking
	return nil
}

func (o *Order) MarkReady() error {
	if o.status != StatusCooking {
		return fmt.Errorf("%w: order is not cooking", ErrInvalidTransition)
	}
	o.status = StatusReady
	return nil
}

func (o *Order) ShipToDelivery() error {
	if o.status != StatusReady {
		return fmt.Errorf("%w: order is not ready", ErrInvalidTransition)
	}
	o.status = StatusDelivering
	return nil
}

func (o *Order) CompleteDelivery() error {
	if o.status != StatusDelivering {
		return fmt.Errorf("%w: order is not in delivery", ErrInvalidTransition)
	}
	o.status = StatusCompleted
	return nil
}

// --- Getters ---

// Items возвращает КОПИЮ позиций.
func (o *Order) Items() []*OrderItem {
	result := make([]*OrderItem, len(o.items))
	// Копируем указатели, но сами OrderItem иммутабельны (у них только геттеры), так что это безопасно.
	copy(result, o.items)
	return result
}
func (o *Order) Address() DeliveryAddress    { return o.address }
func (o *Order) DeliveryPrice() common.Money { return o.deliveryPrice }
func (o *Order) Discount() common.Money      { return o.discount }
func (o *Order) FinalPrice() common.Money    { return o.finalPrice }

// --- Helpers ---

func generateOrderNumber() string {
	return fmt.Sprintf("PG-%s-%s", time.Now().Format("2006.01.02"), uuid.New().String()[:4])
}

// --- Repository ---

type OrderRepository interface {
	Save(o *Order) error
	FindByID(id string) (*Order, error)
}
