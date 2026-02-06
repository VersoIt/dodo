package kitchen

import (
	"testing"
	"time"
)

func TestKitchenTicket_Lifecycle(t *testing.T) {
	items := []KitchenItem{
		{ProductID: "p1", Name: "Pizza", Quantity: 1},
	}
	ticket := NewTicket("order-123", items)

	if ticket.Status() != TicketQueued {
		t.Errorf("expected Queued status")
	}

	_ = ticket.StartCooking()
	if ticket.Status() != TicketCooking {
		t.Errorf("expected Cooking status")
	}

	time.Sleep(10 * time.Millisecond)

	_ = ticket.MarkReady()
	if ticket.Status() != TicketReady {
		t.Errorf("expected Ready status")
	}

	duration := ticket.GetCookingDuration()
	if duration <= 0 {
		t.Errorf("duration should be positive")
	}
}
