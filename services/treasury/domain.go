package treasury

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/versoit/diploma/pkg/common"
)

var (
	ErrPaymentProcessed = errors.New("payment is already processed")
	ErrInvalidRefund    = errors.New("can only refund successful payments")
)

type PaymentMethod int

const (
	MethodOnline PaymentMethod = 0
	MethodCash   PaymentMethod = 1
	MethodCard   PaymentMethod = 2
)

func (m PaymentMethod) String() string {
	switch m {
	case MethodOnline:
		return "online"
	case MethodCash:
		return "cash"
	case MethodCard:
		return "card"
	default:
		return "unknown"
	}
}

type PaymentStatus int

const (
	PayStatusWaiting  PaymentStatus = 0
	PayStatusSuccess  PaymentStatus = 1
	PayStatusDeclined PaymentStatus = 2
	PayStatusRefund   PaymentStatus = 3
)

func (s PaymentStatus) String() string {
	switch s {
	case PayStatusWaiting:
		return "waiting"
	case PayStatusSuccess:
		return "success"
	case PayStatusDeclined:
		return "declined"
	case PayStatusRefund:
		return "refund"
	default:
		return "unknown"
	}
}

type Payment struct {
	id            string
	orderID       string
	amount        common.Money
	method        PaymentMethod
	status        PaymentStatus
	transactionID string
	createdAt     time.Time
	updatedAt     time.Time
}

func NewPayment(orderID string, amount common.Money, method PaymentMethod) *Payment {
	id, _ := uuid.NewV7()
	return &Payment{
		id:        id.String(),
		orderID:   orderID,
		amount:    amount,
		method:    method,
		status:    PayStatusWaiting,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

func (p *Payment) Confirm(externalTransactionID string) error {
	if p.status != PayStatusWaiting {
		return ErrPaymentProcessed
	}
	p.transactionID = externalTransactionID
	p.status = PayStatusSuccess
	p.updatedAt = time.Now()
	return nil
}

func (p *Payment) Decline() error {
	if p.status != PayStatusWaiting {
		return ErrPaymentProcessed
	}
	p.status = PayStatusDeclined
	p.updatedAt = time.Now()
	return nil
}

func (p *Payment) Refund() error {
	if p.status != PayStatusSuccess {
		return ErrInvalidRefund
	}
	p.status = PayStatusRefund
	p.updatedAt = time.Now()
	return nil
}

func (p *Payment) ID() string            { return p.id }
func (p *Payment) OrderID() string       { return p.orderID }
func (p *Payment) Amount() common.Money  { return p.amount }
func (p *Payment) Status() PaymentStatus { return p.status }
func (p *Payment) CreatedAt() time.Time  { return p.createdAt }

type PaymentRepository interface {
	Save(ctx context.Context, p *Payment) error
	FindByOrderID(ctx context.Context, orderID string) (*Payment, error)
}
