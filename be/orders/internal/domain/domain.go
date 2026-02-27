package domain

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
	Entrance  string
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

func ParseStatus(s string) (OrderStatus, error) {
	switch s {
	case "created":
		return StatusCreated, nil
	case "paid":
		return StatusPaid, nil
	case "cooking":
		return StatusCooking, nil
	case "ready":
		return StatusReady, nil
	case "delivering":
		return StatusDelivering, nil
	case "completed":
		return StatusCompleted, nil
	case "canceled":
		return StatusCanceled, nil
	default:
		return -1, fmt.Errorf("unknown status: %s", s)
	}
}

// --- Entities ---

type OrderItem struct {
	id             string
	productID      string
	productName    string
	quantity       Quantity
	basePrice      common.Money
	sizeMultiplier float64
	toppings       []Topping
}

func (i *OrderItem) ID() string {
	if i.id == "" {
		id, _ := uuid.NewV7()
		i.id = id.String()
	}
	return i.id
}

func (i *OrderItem) CalculateTotal() common.Money {
	sizedPrice := i.basePrice.Mul(decimal.NewFromFloat(i.sizeMultiplier))

	toppingsPrice := common.ZeroMoney()
	for _, t := range i.toppings {
		toppingsPrice = toppingsPrice.Add(t.Price)
	}

	unitPrice := sizedPrice.Add(toppingsPrice)
	return unitPrice.Mul(decimal.NewFromInt(int64(i.quantity.Int())))
}

func (i *OrderItem) ProductID() string       { return i.productID }
func (i *OrderItem) ProductName() string     { return i.productName }
func (i *OrderItem) Quantity() int           { return i.quantity.Int() }
func (i *OrderItem) BasePrice() common.Money { return i.basePrice }
func (i *OrderItem) Size() float64           { return i.sizeMultiplier }
func (i *OrderItem) Toppings() []Topping     { return i.toppings }

func (i *OrderItem) AddReconstructedTopping(name string, price common.Money) {
	i.toppings = append(i.toppings, Topping{Name: name, Price: price})
}

// --- Aggregate Root ---

type PromoCode struct {
	id             string
	code           string
	discountType   string // "percent", "fixed"
	discountAmount common.Money
	isActive       bool
	expiresAt      time.Time
}

func (p *PromoCode) CalculateDiscount(basePrice common.Money) common.Money {
	if !p.isActive || (!p.expiresAt.IsZero() && time.Now().After(p.expiresAt)) {
		return common.ZeroMoney()
	}

	if p.discountType == "percent" {
		discount := basePrice.InexactFloat64() * (p.discountAmount.InexactFloat64() / 100.0)
		return common.NewMoney(discount)
	}
	return p.discountAmount
}

func (p *PromoCode) ID() string                   { return p.id }
func (p *PromoCode) Code() string                 { return p.code }
func (p *PromoCode) DiscountType() string         { return p.discountType }
func (p *PromoCode) DiscountAmount() common.Money { return p.discountAmount }
func (p *PromoCode) IsActive() bool               { return p.isActive }

func NewPromoCode(id, code, dType string, amount common.Money, active bool, expires time.Time) *PromoCode {
	return &PromoCode{id, code, dType, amount, active, expires}
}

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

	chefID    string
	courierID string
}

// --- Factory ---

func NewOrder(customerID string, address DeliveryAddress) *Order {
	id, _ := uuid.NewV7()
	return &Order{
		id:            id.String(),
		orderNumber:   generateOrderNumber(),
		customerID:    customerID,
		status:        StatusCreated,
		createdAt:     time.Now(),
		address:       address,
		items:         make([]*OrderItem, 0),
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
	ErrAlreadyAssigned   = errors.New("order is already assigned to someone else")
)

// --- Business Logic ---

func (o *Order) AddItem(productID, name string, qty Quantity, productBasePrice common.Money, sizeMult float64, toppings []Topping) error {
	if o.status != StatusCreated {
		return ErrOrderLocked
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

func (o *Order) SendToKitchen(chefID string) error {
	if o.status != StatusPaid {
		return fmt.Errorf("%w: order must be paid", ErrInvalidTransition)
	}
	if o.chefID != "" && o.chefID != chefID {
		return ErrAlreadyAssigned
	}
	o.status = StatusCooking
	o.chefID = chefID
	return nil
}

func (o *Order) MarkReady() error {
	if o.status != StatusCooking {
		return fmt.Errorf("%w: order is not cooking", ErrInvalidTransition)
	}
	o.status = StatusReady
	return nil
}

func (o *Order) ShipToDelivery(courierID string) error {
	if o.status != StatusReady {
		return fmt.Errorf("%w: order is not ready", ErrInvalidTransition)
	}
	if o.courierID != "" && o.courierID != courierID {
		return ErrAlreadyAssigned
	}
	o.status = StatusDelivering
	o.courierID = courierID
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
func (o *Order) PromoCode() string           { return o.promoCode }
func (o *Order) FinalPrice() common.Money    { return o.finalPrice }
func (o *Order) ChefID() string              { return o.chefID }
func (o *Order) CourierID() string           { return o.courierID }

// ReconstructOrder builds an Order aggregate from storage.
func ReconstructOrder(id, number, custID string, status OrderStatus, createdAt time.Time, addr DeliveryAddress, delPrice, discount common.Money, promo string, final common.Money, items []*OrderItem, chefID, courierID string) *Order {
	return &Order{
		id:            id,
		orderNumber:   number,
		customerID:    custID,
		status:        status,
		createdAt:     createdAt,
		address:       addr,
		items:         items,
		deliveryPrice: delPrice,
		discount:      discount,
		promoCode:     promo,
		finalPrice:    final,
		chefID:        chefID,
		courierID:     courierID,
	}
}

// ReconstructOrderItem builds an OrderItem from storage.
func ReconstructOrderItem(id, prodID, name string, qty int, base common.Money, size float64, toppings []Topping) *OrderItem {
	// Trust DB data for reconstruction
	q, _ := NewQuantity(qty)
	return &OrderItem{
		id:             id,
		productID:      prodID,
		productName:    name,
		quantity:       q,
		basePrice:      base,
		sizeMultiplier: size,
		toppings:       toppings,
	}
}

func generateOrderNumber() string {
	id, _ := uuid.NewV7()
	// Используем последние 8 символов UUID для большей уникальности
	return fmt.Sprintf("PG-%s-%s", time.Now().Format("2006.01.02"), id.String()[24:])
}

type OrderFilter struct {
	StartAt *time.Time
	EndAt   *time.Time
	Status  *OrderStatus
}

type OrderRepository interface {
	Save(ctx context.Context, o *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*Order, error)
	FindAll(ctx context.Context) ([]*Order, error)
	FindFiltered(ctx context.Context, filter OrderFilter) ([]*Order, error)

	// Promo codes
	SavePromo(ctx context.Context, p *PromoCode) error
	FindPromoByCode(ctx context.Context, code string) (*PromoCode, error)
	ListPromos(ctx context.Context) ([]*PromoCode, error)
	DeletePromo(ctx context.Context, id string) error
}

type ProductInfo struct {
	ID        string
	Name      string
	BasePrice common.Money
}

type CatalogService interface {
	GetProduct(ctx context.Context, id string) (*ProductInfo, error)
}

type KitchenService interface {
	CreateTicket(ctx context.Context, orderID string, orderNumber string, items []*OrderItem) error
}

type LogisticsService interface {
	CreateDelivery(ctx context.Context, orderID string, orderNumber string, address DeliveryAddress, items []*OrderItem) error
}

type TreasuryService interface {
	ProcessPayment(ctx context.Context, orderID string, amount common.Money) error
}

type NotificationService interface {
	NotifyStatusChanged(ctx context.Context, customerID string, orderID string, status OrderStatus) error
}
