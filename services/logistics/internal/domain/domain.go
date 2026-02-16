package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type DeliveryStatus int

const (
	DelStatusPending   DeliveryStatus = 0
	DelStatusAssigned  DeliveryStatus = 1
	DelStatusOnWay     DeliveryStatus = 2
	DelStatusDelivered DeliveryStatus = 3
	DelStatusFailed    DeliveryStatus = 4
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

type Delivery struct {
	orderID      string
	orderNumber  string
	courierID    string
	status       DeliveryStatus
	createdAt    time.Time
	pickupTime   time.Time
	deliveryTime time.Time

	city      string
	street    string
	house     string
	apartment string
	items     []DeliveryItem

	location Coordinates
}

type DeliveryItem struct {
	ProductID string
	Name      string
	Quantity  int
}

type Courier struct {
	id       string
	name     string
	phone    string
	status   CourierStatus
	location Coordinates
}

func NewDelivery(orderID, orderNumber, city, street, house, apartment string, items []DeliveryItem) *Delivery {
	return &Delivery{
		orderID:     orderID,
		orderNumber: orderNumber,
		city:        city,
		street:      street,
		house:       house,
		apartment:   apartment,
		items:       items,
		status:      DelStatusPending,
		createdAt:   time.Now(),
	}
}

func NewCourier(name, phone string) *Courier {
	id, _ := uuid.NewV7()
	return &Courier{
		id:     id.String(),
		name:   name,
		phone:  phone,
		status: CourierOffline,
	}
}

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

var (
	ErrDeliveryNotPending = errors.New("delivery is not in pending state")
	ErrCourierNotAssigned = errors.New("courier is not assigned")
	ErrInvalidStatus      = errors.New("invalid status for operation")
	ErrCourierBusy        = errors.New("courier is busy")
	ErrInvalidCoordinates = errors.New("invalid coordinates")
)

func (d *Delivery) UpdateLocation(coords Coordinates) {
	d.location = coords
}

func (c *Courier) UpdateLocation(coords Coordinates) {
	c.location = coords
}

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

func (d *Delivery) ID() string {
	// For Delivery, orderID is effectively the unique ID if 1 order = 1 delivery
	return d.orderID
}

func (d *Delivery) OrderID() string              { return d.orderID }
func (d *Delivery) OrderNumber() string          { return d.orderNumber }
func (d *Delivery) CourierID() string            { return d.courierID }
func (d *Delivery) Status() DeliveryStatus       { return d.status }
func (d *Delivery) CreatedAt() time.Time         { return d.createdAt }
func (d *Delivery) PickupTime() time.Time        { return d.pickupTime }
func (d *Delivery) DeliveryTime() time.Time      { return d.deliveryTime }
func (d *Delivery) Location() (lat, lng float64) { return d.location.Lat, d.location.Lng }
func (d *Delivery) City() string                 { return d.city }
func (d *Delivery) Street() string               { return d.street }
func (d *Delivery) House() string                { return d.house }
func (d *Delivery) Apartment() string            { return d.apartment }
func (d *Delivery) Items() []DeliveryItem        { return d.items }

func (c *Courier) ID() string                   { return c.id }
func (c *Courier) Name() string                 { return c.name }
func (c *Courier) Phone() string                { return c.phone }
func (c *Courier) Status() CourierStatus        { return c.status }
func (c *Courier) Location() (lat, lng float64) { return c.location.Lat, c.location.Lng }

// ReconstructCourier builds a Courier from storage.
func ReconstructCourier(id, name, phone string, status CourierStatus, lat, lng float64) *Courier {
	return &Courier{
		id:       id,
		name:     name,
		phone:    phone,
		status:   status,
		location: Coordinates{Lat: lat, Lng: lng},
	}
}

// ReconstructDelivery builds a Delivery from storage.
func ReconstructDelivery(orderID, orderNumber, courierID string, status DeliveryStatus, createdAt, pickup, del time.Time, lat, lng float64, city, street, house, apartment string, items []DeliveryItem) *Delivery {
	return &Delivery{
		orderID:      orderID,
		orderNumber:  orderNumber,
		courierID:    courierID,
		status:       status,
		createdAt:    createdAt,
		pickupTime:   pickup,
		deliveryTime: del,
		location:     Coordinates{Lat: lat, Lng: lng},
		city:         city,
		street:       street,
		house:        house,
		apartment:    apartment,
		items:        items,
	}
}

type DeliveryRepository interface {
	Save(ctx context.Context, d *Delivery) error
	FindByOrderID(ctx context.Context, orderID string) (*Delivery, error)
	FindAll(ctx context.Context) ([]*Delivery, error)
}

type CourierRepository interface {
	FindAvailable(ctx context.Context) ([]*Courier, error)
	FindByID(ctx context.Context, id string) (*Courier, error)
	Save(ctx context.Context, c *Courier) error
	UpdateLocation(ctx context.Context, id string, lat, lng float64) error
}
