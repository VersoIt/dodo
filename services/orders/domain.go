package orders

import (
	"diploma/pkg/common"
	"errors"
	"fmt"
	"time"

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
	
	items       []*OrderItem
	address     DeliveryAddress
	
	deliveryPrice common.Money
	discount      common.Money
	promoCode     string
	
	finalPrice    common.Money 
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

// --- Business Logic ---

// AddItem добавляет товар.
// Передаем basePrice (чистую) и sizeMultiplier отдельно.
func (o *Order) AddItem(productID, name string, qty int, basePrice common.Money, sizeMult float64, toppings []Topping) error {
	if o.status != StatusCreated {
		return errors.New("cannot add items to processed order")
	}
	if qty <= 0 {
		return errors.New("quantity must be positive")
	}
	if sizeMult <= 0 {
		sizeMult = 1.0
	}

	item := &OrderItem{
		productID:      productID,
		productName:    name,
		quantity:       qty,
		basePrice:      basePrice,
		sizeMultiplier: sizeMult,
		toppings:       toppings,
	}
	
	o.items = append(o.items, item)
	o.recalculate()
	return nil
}

func (o *Order) ApplyPromoCode(code string, discountAmount common.Money) error {
	if o.status != StatusCreated {
		return errors.New("cannot apply promo to processed order")
	}
	o.promoCode = code
	o.discount = discountAmount
	o.recalculate()
	return nil
}

func (o *Order) SetDeliveryPrice(price common.Money) {
	if o.status != StatusCreated {
		return 
	}
	o.deliveryPrice = price
	o.recalculate()
}

func (o *Order) recalculate() {
	var itemsTotal common.Money
	for _, item := range o.items {
		itemsTotal += item.CalculateTotal()
	}
	
	total := itemsTotal + o.deliveryPrice - o.discount
	if total < 0 {
		total = 0
	}
	o.finalPrice = total
}

// --- State Machine ---

func (o *Order) MarkPaid() error {
	if o.status != StatusCreated {
		return fmt.Errorf("invalid transition: Created -> Paid from status %v", o.status)
	}
	o.status = StatusPaid
	return nil
}

func (o *Order) SendToKitchen() error {
	if o.status != StatusPaid {
		return errors.New("order must be paid before cooking")
	}
	o.status = StatusCooking
	return nil
}

func (o *Order) MarkReady() error {
	if o.status != StatusCooking {
		return errors.New("order is not cooking")
	}
	o.status = StatusReady
	return nil
}

func (o *Order) ShipToDelivery() error {
	if o.status != StatusReady {
		return errors.New("order is not ready for delivery")
	}
	o.status = StatusDelivering
	return nil
}

func (o *Order) CompleteDelivery() error {
	if o.status != StatusDelivering {
		return errors.New("order was not being delivered")
	}
	o.status = StatusCompleted
	return nil
}

// --- Getters ---

func (o *Order) ID() string                  { return o.id }
func (o *Order) OrderNumber() string         { return o.orderNumber }
func (o *Order) CustomerID() string          { return o.customerID }
func (o *Order) Status() OrderStatus         { return o.status }
func (o *Order) CreatedAt() time.Time        { return o.createdAt }
func (o *Order) Items() []*OrderItem         { return o.items }
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
