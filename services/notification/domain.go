package notification

import (
	"time"

	"github.com/google/uuid"
)

// --- Enums ---

type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelPush  Channel = "push"
)

// --- Aggregate ---

type Notification struct {
	id      string
	userID  string
	channel Channel
	title   string
	message string
	sentAt  time.Time
	isSent  bool
	error   string
}

// --- Factory ---

func NewNotification(userID string, ch Channel, title, msg string) *Notification {
	return &Notification{
		id:      uuid.New().String(),
		userID:  userID,
		channel: ch,
		title:   title,
		message: msg,
	}
}

// --- Behavior ---

func (n *Notification) MarkSent() {
	n.isSent = true
	n.sentAt = time.Now()
}

func (n *Notification) MarkFailed(errStr string) {
	n.isSent = false
	n.error = errStr
}

// Getters
func (n *Notification) ID() string { return n.id }
func (n *Notification) UserID() string { return n.userID }
func (n *Notification) Title() string { return n.title }
func (n *Notification) Message() string { return n.message }
func (n *Notification) SentAt() time.Time { return n.sentAt }
func (n *Notification) IsSent() bool { return n.isSent }
func (n *Notification) Error() string { return n.error }

// --- Repository ---

type NotificationRepository interface {
	Save(n *Notification) error
}
