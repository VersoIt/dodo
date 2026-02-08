package logistics

import (
	"testing"
)

func TestCourier_Workflow(t *testing.T) {
	c := NewCourier("John", "123456")

	c.GoOnline()
	if c.Status() != CourierFree {
		t.Errorf("courier should be free")
	}

	_ = c.TakeOrder()
	if c.Status() != CourierBusy {
		t.Errorf("courier should be busy")
	}

	err := c.GoOffline()
	if err != ErrCourierBusy {
		t.Errorf("busy courier cannot go offline")
	}
}

func TestDelivery_Lifecycle(t *testing.T) {
	d := NewDelivery("order-1")

	_ = d.AssignCourier("c-1")
	if d.Status() != DelStatusAssigned {
		t.Errorf("status mismatch")
	}

	_ = d.Pickup()
	_ = d.Complete()

	if d.Status() != DelStatusDelivered {
		t.Errorf("delivery should be completed")
	}
}
