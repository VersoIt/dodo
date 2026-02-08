package orders

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/versoit/diploma/pkg/common"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

type OrderItem struct {
	productID      string
	productName    string
	quantity       int
	basePrice      common.Money
	sizeMultiplier float64
	toppings       []Topping
}

func (i *OrderItem) CalculateTotal() common.Money {
	sizedPrice := i.basePrice.Mul(decimal.NewFromFloat(i.sizeMultiplier))

	toppingsPrice := common.ZeroMoney()
	for _, t := range i.toppings {
		toppingsPrice = toppingsPrice.Add(t.Price)
	}

	unitPrice := sizedPrice.Add(toppingsPrice)
	return unitPrice.Mul(decimal.NewFromInt(int64(i.quantity)))
}

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
	id, _ := uuid.NewV7()
	return &Order{
		id:          id.String(),
		orderNumber: generateOrderNumber(),
		customerID:  customerID,
		status:      StatusCreated,
		createdAt:   time.Now(),
		address:     address,
		items:       make([]*OrderItem, 0),
		deliveryPrice: common.ZeroMoney(),
		discount:      common.ZeroMoney(),
		finalPrice:    common.ZeroMoney(),
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

// --- Business Logic ---

func (o *Order) AddItem(productID, name string, qty int, productBasePrice common.Money, sizeMult float64, toppings []Topping) error {
	if o.status != StatusCreated {
		return ErrOrderLocked
	}
	if qty <= 0 {
		return ErrInvalidQty
	}

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
	if discountAmount.IsNegative() {
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
	total := common.ZeroMoney()
	for _, item := range o.items {
		total = total.Add(item.CalculateTotal())
	}

	o.finalPrice = total.Add(o.deliveryPrice).Sub(o.discount)
	if o.finalPrice.IsNegative() {
		o.finalPrice = common.ZeroMoney()
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

func (o *Order) ID() string           { return o.id }
func (o *Order) OrderNumber() string  { return o.orderNumber }
func (o *Order) CustomerID() string   { return o.customerID }
func (o *Order) Status() OrderStatus  { return o.status }
func (o *Order) CreatedAt() time.Time { return o.createdAt }
func (o *Order) Items() []*OrderItem {
	result := make([]*OrderItem, len(o.items))
	copy(result, o.items)
	return result
}
func (o *Order) Address() DeliveryAddress    { return o.address }
func (o *Order) DeliveryPrice() common.Money { return o.deliveryPrice }
func (o *Order) Discount() common.Money      { return o.discount }
func (o *Order) FinalPrice() common.Money    { return o.finalPrice }

func generateOrderNumber() string {
	id, _ := uuid.NewV7()
	return fmt.Sprintf("PG-%s-%s", time.Now().Format("2006.01.02"), id.String()[:4])
}

type OrderRepository interface {
	Save(ctx context.Context, o *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
}
