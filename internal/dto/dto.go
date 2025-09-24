package dto

import (
	"time"
)

type BookingDTO struct {
	ID         int       `json:"id"`
	EventID    int       `json:"event_id"`
	EventName  string    `json:"event_name"`
	TelegramID int       `json:"telegram_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type QueueMessage struct {
	BookingID   int `json:"booking_id"`
	PlacesCount int `json:"places_count"`
}

type CreateEvent struct {
	EventName string    `json:"event_name"`
	EventAt   time.Time `json:"event_at"`
	AllSeats  int       `json:"all_seats"`
}

type CreateBooking struct {
	EventID     int `json:"event_id,omitempty"`
	TelegramID  int `json:"telegram_id"`
	PlacesCount int `json:"places_count"`
}
