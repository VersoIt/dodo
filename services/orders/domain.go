package orders

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// --- Value Objects ---

type Money float64

type OrderStatus int

const (
	StatusCreated    OrderStatus = 0
	StatusPaid       OrderStatus = 1 // Оплачен, ждет кухни
	StatusCooking    OrderStatus = 2
	StatusReady      OrderStatus = 3
	StatusDelivering OrderStatus = 4
	StatusCompleted  OrderStatus = 5
	StatusCanceled   OrderStatus = 6
)

// DeliveryAddress - Value Object. Неизменяем.
type DeliveryAddress struct {
	City      string
	Street    string
	House     string
	Apartment string
	Floor     string
	Comment   string
}

// Topping - добавка. Value Object.
type Topping struct {
	Name  string
	Price Money
}

// --- Entities ---

// OrderItem - Позиция заказа.
// Считает свою стоимость сама.
type OrderItem struct {
	productID      string
	productName    string
	quantity       int
	basePrice      Money // Цена за 1 шт (уже с учетом size_m)
	sizeMultiplier float64
	toppings       []Topping
}

// CalculateTotal возвращает (Base * Qty) + ToppingsCost.
// Формула из ТЗ (1): S_ki = (base_i * size_m + SUM(top_ij)) * (1 - d_k)
// Здесь мы считаем часть до скидки: (base * size + toppings) * qty
func (i *OrderItem) CalculateTotal() Money {
	itemPrice := i.basePrice // basePrice уже может включать size_m, если так передали, или умножим тут
	// В CatalogService мы сделали CalculateBasePrice(size).
	// Допустим, сюда приходит уже "цена за размер".

	var toppingsPrice Money
	for _, t := range i.toppings {
		toppingsPrice += t.Price
	}

	// Цена за единицу = (База + Топпинги)
	unitPrice := itemPrice + toppingsPrice

	return unitPrice * Money(i.quantity)
}

// --- Aggregate Root ---

// Order - Агрегат заказа.
// Гарантирует консистентность состояния и расчетов.
type Order struct {
	id          string
	orderNumber string // PG-YYYY.MM.DD-NNNN
	customerID  string
	status      OrderStatus
	createdAt   time.Time
	
	items       []*OrderItem
	address     DeliveryAddress
	
	// Финансы
	deliveryPrice Money
	discount      Money // Скидка в деньгах
	promoCode     string
	
	// Кеш итоговой суммы (пересчитывается при изменении)
	finalPrice    Money 
}

// --- Factory ---

func NewOrder(customerID string, address DeliveryAddress) *Order {
	return &Order{
		id:          uuid.New().String(),
		orderNumber: generateOrderNumber(), // Заглушка, реальный генератор сложнее
		customerID:  customerID,
		status:      StatusCreated,
		createdAt:   time.Now(),
		address:     address,
		items:       make([]*OrderItem, 0),
	}
}

// --- Business Logic ---

// AddItem добавляет товар.
// basePriceWithSize - это цена товара с учетом размера (из Catalog Service).
func (o *Order) AddItem(productID, name string, qty int, basePriceWithSize Money, toppings []Topping) error {
	if o.status != StatusCreated {
		return errors.New("cannot add items to processed order")
	}
	if qty <= 0 {
		return errors.New("quantity must be positive")
	}

	item := &OrderItem{
		productID:   productID,
		productName: name,
		quantity:    qty,
		basePrice:   basePriceWithSize,
		toppings:    toppings,
	}
	
	o.items = append(o.items, item)
	o.recalculate()
	return nil
}

// ApplyPromoCode применяет скидку.
// В реальности здесь была бы сложная логика проверки промокода.
func (o *Order) ApplyPromoCode(code string, discountAmount Money) error {
	if o.status != StatusCreated {
		return errors.New("cannot apply promo to processed order")
	}
	o.promoCode = code
	o.discount = discountAmount
	o.recalculate()
	return nil
}

// SetDeliveryPrice устанавливает цену доставки.
func (o *Order) SetDeliveryPrice(price Money) {
	if o.status != StatusCreated {
		return 
	}
	o.deliveryPrice = price
	o.recalculate()
}

// recalculate - внутренняя метод для обновления суммы.
// Формула (2): S_total = SUM(S_ki) + Delivery
func (o *Order) recalculate() {
	var itemsTotal Money
	for _, item := range o.items {
		itemsTotal += item.CalculateTotal()
	}
	
	total := itemsTotal + o.deliveryPrice - o.discount
	if total < 0 {
		total = 0
	}
	o.finalPrice = total
}

// --- State Machine Transitions ---

// MarkPaid переводит в статус "Оплачен".
func (o *Order) MarkPaid() error {
	if o.status != StatusCreated {
		return fmt.Errorf("invalid transition: Created -> Paid from status %v", o.status)
	}
	o.status = StatusPaid
	return nil
}

// SendToKitchen отправляет на кухню.
func (o *Order) SendToKitchen() error {
	if o.status != StatusPaid {
		return errors.New("order must be paid before cooking")
	}
	o.status = StatusCooking
	return nil
}

// MarkReady - кухня закончила.
func (o *Order) MarkReady() error {
	if o.status != StatusCooking {
		return errors.New("order is not cooking")
	}
	o.status = StatusReady
	return nil
}

// ShipToDelivery - отдали курьеру.
func (o *Order) ShipToDelivery() error {
	if o.status != StatusReady {
		return errors.New("order is not ready for delivery")
	}
	o.status = StatusDelivering
	return nil
}

// CompleteDelivery - успешно доставлено.
func (o *Order) CompleteDelivery() error {
	if o.status != StatusDelivering {
		return errors.New("order was not being delivered")
	}
	o.status = StatusCompleted
	return nil
}

// --- Getters ---

func (o *Order) ID() string { return o.id }
func (o *Order) FinalPrice() Money { return o.finalPrice }
func (o *Order) Status() OrderStatus { return o.status }
func (o *Order) Items() []*OrderItem { return o.items }

// --- Helpers ---

func generateOrderNumber() string {
	// В реальности: PG-2026.02.06-0001
	// Для MVP:
	return fmt.Sprintf("PG-%s-%s", time.Now().Format("2006.01.02"), uuid.New().String()[:4])
}

// --- Repository ---

type OrderRepository interface {
	Save(o *Order) error
	FindByID(id string) (*Order, error)
}
