package model

import "time"

type Event struct {
	ID        int       `json:"id"`
	EventName string    `json:"event_name"`
	AllSeats  int       `json:"all_seats"`
	Booked    int       `json:"booked"`
	EventAt   time.Time `json:"event_at"`
	CreatedAt time.Time `json:"created_at"`
	Bookings  []Booking `json:"bookings,omitempty"`
}

type Booking struct {
	ID          int       `json:"id"`
	EventID     int       `json:"event_id"`
	PlacesCount int       `json:"places_count"`
	Status      string    `json:"status"`
	TelegramID  int       `json:"telegram_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
