package treasury

import (
	"diploma/pkg/common"
	"errors"
	"time"

	"github.com/google/uuid"
)

// --- Enums ---

type PaymentMethod int

const (
	MethodOnline PaymentMethod = 0
	MethodCash   PaymentMethod = 1
	MethodCard   PaymentMethod = 2
)

type PaymentStatus int

const (
	PayStatusWaiting  PaymentStatus = 0
	PayStatusSuccess  PaymentStatus = 1
	PayStatusDeclined PaymentStatus = 2
	PayStatusRefund   PaymentStatus = 3
)

// --- Aggregate ---

// Payment - Агрегат платежа.
type Payment struct {
	id            string
	orderID       string
	amount        common.Money
	method        PaymentMethod
	status        PaymentStatus
	transactionID string // ID из банка
	createdAt     time.Time
	updatedAt     time.Time
}

// --- Factory ---

func NewPayment(orderID string, amount common.Money, method PaymentMethod) *Payment {
	return &Payment{
		id:        uuid.New().String(),
		orderID:   orderID,
		amount:    amount,
		method:    method,
		status:    PayStatusWaiting,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

// --- Behavior ---

func (p *Payment) Confirm(externalTransactionID string) error {
	if p.status != PayStatusWaiting {
		return errors.New("payment is already processed")
	}
	p.transactionID = externalTransactionID
	p.status = PayStatusSuccess
	p.updatedAt = time.Now()
	return nil
}

func (p *Payment) Decline() error {
	if p.status != PayStatusWaiting {
		return errors.New("payment is already processed")
	}
	p.status = PayStatusDeclined
	p.updatedAt = time.Now()
	return nil
}

func (p *Payment) Refund() error {
	if p.status != PayStatusSuccess {
		return errors.New("can only refund successful payments")
	}
	p.status = PayStatusRefund
	p.updatedAt = time.Now()
	return nil
}

// --- Getters ---

func (p *Payment) ID() string            { return p.id }
func (p *Payment) OrderID() string       { return p.orderID }
func (p *Payment) Amount() common.Money  { return p.amount }
func (p *Payment) Status() PaymentStatus { return p.status }
func (p *Payment) CreatedAt() time.Time  { return p.createdAt }

// --- Repository ---

type PaymentRepository interface {
	Save(p *Payment) error
	FindByOrderID(orderID string) (*Payment, error)
}
