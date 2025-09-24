package repository

import (
	"errors"
	"fmt"

	"github.com/Komilov31/event-booker/internal/dto"
	"github.com/Komilov31/event-booker/internal/model"
	"github.com/lib/pq"
)

func (r *Repository) CreateEvent(event dto.CreateEvent) (*model.Event, error) {
	query := `INSERT INTO events(event_name, event_at, all_seats) 
	VALUES ($1, $2, $3) RETURNING id, created_at`

	var createdEvent model.Event
	err := r.db.Master.QueryRow(query, event.EventName, event.EventAt, event.AllSeats).Scan(
		&createdEvent.ID,
		&createdEvent.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create event in db: %w", err)
	}

	createdEvent.EventName = event.EventName
	createdEvent.EventAt = event.EventAt
	createdEvent.AllSeats = event.AllSeats

	return &createdEvent, nil
}

func (r *Repository) CreateBooking(booking dto.CreateBooking) (*model.Booking, error) {
	tx, err := r.db.Master.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not start transcation: %w", err)
	}
	defer tx.Rollback()

	var createdBooking model.Booking
	query := `INSERT INTO bookings(event_id, status, telegram_id)
	VALUES ($1, $2, $3) RETURNING id, created_at`
	err = tx.QueryRow(
		query,
		booking.EventID,
		"created",
		booking.TelegramID,
	).Scan(&createdBooking.ID, &createdBooking.CreatedAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, ErrNoSuchEvent
		}

		return nil, fmt.Errorf("could not create booking in db: %w", err)
	}

	query = `UPDATE events
	SET booked = booked + $1
	WHERE id = $2 AND booked + $1 <= all_seats`
	result, err := tx.Exec(query, booking.PlacesCount, booking.EventID)
	if err != nil {
		return nil, fmt.Errorf("could not book in event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("could not book in event: %w", err)
	}

	if rowsAffected == 0 {
		return nil, ErrNoSeatsAvailable
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transcation: %w", err)
	}

	createdBooking.EventID = booking.EventID
	createdBooking.PlacesCount = booking.PlacesCount
	createdBooking.TelegramID = booking.TelegramID
	createdBooking.Status = "created"

	return &createdBooking, nil
}
