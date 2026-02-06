package treasury

import (
	"diploma/pkg/common"
	"testing"
)

func TestPayment_Confirm(t *testing.T) {
	p := NewPayment("order-1", common.Money(100.0), MethodOnline)

	err := p.Confirm("tx-999")
	if err != nil {
		t.Fatalf("confirm failed: %v", err)
	}

	if p.Status() != PayStatusSuccess {
		t.Errorf("expected success status")
	}
}

func TestPayment_Refund(t *testing.T) {
	p := NewPayment("order-1", 100, MethodOnline)
	
	// Cannot refund waiting payment
	if err := p.Refund(); err != ErrInvalidRefund {
		t.Errorf("expected ErrInvalidRefund")
	}

	_ = p.Confirm("tx-1")
	if err := p.Refund(); err != nil {
		t.Errorf("refund failed: %v", err)
	}

	if p.Status() != PayStatusRefund {
		t.Errorf("expected refund status")
	}
}
