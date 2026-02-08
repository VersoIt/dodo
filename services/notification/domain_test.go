package notification

import (
	"testing"
)

func TestNotification_Lifecycle(t *testing.T) {
	n := NewNotification("user1", ChannelEmail, "Subject", "Body")

	if n.IsSent() {
		t.Error("new notification should not be sent")
	}

	n.MarkSent()
	if !n.IsSent() {
		t.Error("notification should be marked as sent")
	}
	if n.SentAt().IsZero() {
		t.Error("sent time should be recorded")
	}

	n.MarkFailed("error msg")
	if n.IsSent() {
		t.Error("failed notification should not be marked as sent (or logic depends on semantics, but let's assume retry)")
	}
	if n.Error() != "error msg" {
		t.Error("error message mismatch")
	}
}
