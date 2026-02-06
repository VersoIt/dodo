package logistics

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// --- Enums ---

type DeliveryStatus int

const (
	DelStatusPending   DeliveryStatus = 0 // Ждет курьера
	DelStatusAssigned  DeliveryStatus = 1 // Курьер назначен
	DelStatusOnWay     DeliveryStatus = 2 // Забрал, везет
	DelStatusDelivered DeliveryStatus = 3 // Успех
	DelStatusFailed    DeliveryStatus = 4 // Неудача
)

type CourierStatus int

const (
	CourierOffline CourierStatus = 0
	CourierFree    CourierStatus = 1
	CourierBusy    CourierStatus = 2
)

func (s DeliveryStatus) String() string {
	switch s {
	case DelStatusPending:
		return "pending"
	case DelStatusAssigned:
		return "assigned"
	case DelStatusOnWay:
		return "on_way"
	case DelStatusDelivered:
		return "delivered"
	case DelStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func (s CourierStatus) String() string {
	switch s {
	case CourierOffline:
		return "offline"
	case CourierFree:
		return "free"
	case CourierBusy:
		return "busy"
	default:
		return "unknown"
	}
}

// --- Aggregates ---

type Delivery struct {
	orderID      string
	courierID    string
	status       DeliveryStatus
	createdAt    time.Time
	pickupTime   time.Time
	deliveryTime time.Time
	
	currentLat float64
	currentLng float64
}

type Courier struct {
	id         string
	name       string
	phone      string
	status     CourierStatus
	currentLat float64
	currentLng float64
}

// --- Factories ---

func NewDelivery(orderID string) *Delivery {
	return &Delivery{
		orderID:   orderID,
		status:    DelStatusPending,
		createdAt: time.Now(),
	}
}

func NewCourier(name, phone string) *Courier {
	return &Courier{
		id:     uuid.New().String(),
		name:   name,
		phone:  phone,
		status: CourierOffline,
	}
}

// --- Behavior: Delivery ---

func (d *Delivery) AssignCourier(courierID string) error {
	if d.status != DelStatusPending {
		return ErrDeliveryNotPending
	}
	d.courierID = courierID
	d.status = DelStatusAssigned
	return nil
}

func (d *Delivery) Pickup() error {
	if d.status != DelStatusAssigned {
		return ErrCourierNotAssigned
	}
	d.status = DelStatusOnWay
	d.pickupTime = time.Now()
	return nil
}

func (d *Delivery) Complete() error {
	if d.status != DelStatusOnWay {
		return ErrInvalidStatus
	}
	d.status = DelStatusDelivered
	d.deliveryTime = time.Now()
	return nil
}

// --- Errors ---

var (
	ErrDeliveryNotPending  = errors.New("delivery is not in pending state")
	ErrCourierNotAssigned  = errors.New("courier is not assigned")
	ErrInvalidStatus       = errors.New("invalid status for operation")
	ErrCourierBusy         = errors.New("courier is busy")
	ErrInvalidCoordinates  = errors.New("invalid coordinates")
)

// ...

func (d *Delivery) UpdateLocation(lat, lng float64) error {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return ErrInvalidCoordinates
	}
	d.currentLat = lat
	d.currentLng = lng
	return nil
}

// ...

func (c *Courier) UpdateLocation(lat, lng float64) error {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return ErrInvalidCoordinates
	}
	c.currentLat = lat
	c.currentLng = lng
	return nil
}

// --- Behavior: Courier ---

func (c *Courier) GoOnline() {
	if c.status == CourierOffline {
		c.status = CourierFree
	}
}

func (c *Courier) GoOffline() error {
	if c.status == CourierBusy {
		return ErrCourierBusy
	}
	c.status = CourierOffline
	return nil
}

func (c *Courier) TakeOrder() error {
	if c.status != CourierFree {
		return ErrCourierBusy
	}
	c.status = CourierBusy
	return nil
}

func (c *Courier) CompleteOrder() {
	if c.status == CourierBusy {
		c.status = CourierFree
	}
}

// Getters
func (d *Delivery) OrderID() string { return d.orderID }
func (d *Delivery) CourierID() string { return d.courierID }
func (d *Delivery) Status() DeliveryStatus { return d.status }
func (d *Delivery) PickupTime() time.Time { return d.pickupTime }
func (d *Delivery) DeliveryTime() time.Time { return d.deliveryTime }
func (d *Delivery) Location() (lat, lng float64) { return d.currentLat, d.currentLng }

func (c *Courier) ID() string { return c.id }
func (c *Courier) Name() string { return c.name }
func (c *Courier) Phone() string { return c.phone }
func (c *Courier) Status() CourierStatus { return c.status }
func (c *Courier) Location() (lat, lng float64) { return c.currentLat, c.currentLng }

// --- Repository ---

type DeliveryRepository interface {
	Save(d *Delivery) error
	FindByOrderID(orderID string) (*Delivery, error)
}

type CourierRepository interface {
	FindAvailable() ([]*Courier, error)
	Save(c *Courier) error
	UpdateLocation(id string, lat, lng float64) error
}
