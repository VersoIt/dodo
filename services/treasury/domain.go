package treasury

import (
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

// Money - дублируем тип, чтобы сервис был независим (Bounded Context).
type Money float64

// --- Aggregate ---

// Payment - Агрегат платежа.
type Payment struct {
	id            string
	orderID       string
	amount        Money
	method        PaymentMethod
	status        PaymentStatus
	transactionID string // ID из банка
	createdAt     time.Time
	updatedAt     time.Time
}

// --- Factory ---

func NewPayment(orderID string, amount Money, method PaymentMethod) *Payment {
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

// Confirm подтверждает платеж (успех).
func (p *Payment) Confirm(externalTransactionID string) error {
	if p.status != PayStatusWaiting {
		return errors.New("payment is already processed")
	}
	p.transactionID = externalTransactionID
	p.status = PayStatusSuccess
	p.updatedAt = time.Now()
	return nil
}

// Decline отклоняет платеж.
func (p *Payment) Decline() error {
	if p.status != PayStatusWaiting {
		return errors.New("payment is already processed")
	}
	p.status = PayStatusDeclined
	p.updatedAt = time.Now()
	return nil
}

// Refund делает возврат.
func (p *Payment) Refund() error {
	if p.status != PayStatusSuccess {
		return errors.New("can only refund successful payments")
	}
	p.status = PayStatusRefund
	p.updatedAt = time.Now()
	return nil
}

// Getters
func (p *Payment) ID() string { return p.id }
func (p *Payment) Status() PaymentStatus { return p.status }

// --- Repository ---

type PaymentRepository interface {
	Save(p *Payment) error
	FindByOrderID(orderID string) (*Payment, error)
}
