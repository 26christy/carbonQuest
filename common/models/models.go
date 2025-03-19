package models

import (
	"time"

	"github.com/google/uuid"
)

type Alarm struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" validate:"required,min=3,max=100"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateAlarm struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status" validate:"oneof=triggered active ACK"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AlarmEvent struct {
	AlarmID   string    `json:"alarm_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Type      string    `json:"type" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
}

type ACKState struct {
	AlarmID            string    `json:"alarm_id"`
	ACKedAt            time.Time `json:"acked_at"`
	NextNotificationAt time.Time `json:"next_notification_at"`
	ShouldNotify       bool      `json:"should_notify"`
}
