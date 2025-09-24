package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Komilov31/event-booker/internal/dto"
	"github.com/Komilov31/event-booker/internal/model"
)

func (r *Repository) GetEventByID(id int) (*model.Event, error) {
	query := "SELECT * FROM events WHERE id = $1"

	var event model.Event
	err := r.db.Master.QueryRow(query, id).Scan(
		&event.ID,
		&event.EventName,
		&event.AllSeats,
		&event.Booked,
		&event.EventAt,
		&event.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoSuchEvent
		}
		return nil, fmt.Errorf("could not get event from db: %w", err)
	}

	return &event, nil
}

func (r *Repository) GetBookingByID(id int) (*dto.BookingDTO, error) {
	query := `SELECT b.*, e.event_name 
	FROM bookings b 
	JOIN events e ON e.id = b.event_id
	WHERE b.id = $1`

	var booking dto.BookingDTO
	err := r.db.Master.QueryRow(query, id).Scan(
		&booking.ID,
		&booking.EventID,
		&booking.Status,
		&booking.TelegramID,
		&booking.CreatedAt,
		&booking.EventName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoSuchBooking
		}

		return nil, fmt.Errorf("could not get booking info from db: %w", err)
	}

	return &booking, nil
}

func (r *Repository) GetEventWithBookingsByID(id int) (*model.Event, error) {
	event, err := r.GetEventByID(id)
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM bookings WHERE event_id = $1"

	rows, err := r.db.Master.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("could not get bookings from db: %w", err)
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.EventID,
			&booking.Status,
			&booking.TelegramID,
			&booking.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not get bookings from db: %w", err)
		}

		bookings = append(bookings, booking)
	}

	event.Bookings = bookings
	return event, nil
}

func (r *Repository) GetAllEvents() ([]model.Event, error) {
	query := `SELECT 
  	e.id,
  	e.event_name,
  	e.all_seats,
  	e.booked,
  	e.created_at,
  	COALESCE(json_agg(json_build_object(
   	 	'id', b.id,
    	'event_id', b.event_id,
    	'status', b.status,
    	'telegram_id', b.telegram_id,
    	'created_at', b.created_at
  	)) FILTER (WHERE b.id IS NOT NULL), '[]') AS bookings
	FROM events e
	LEFT JOIN bookings b ON b.event_id = e.id
	GROUP BY e.id;`

	rows, err := r.db.Master.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get events from db: %w", err)
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var bookingsJson []byte
		var event model.Event
		var bookings []model.Booking
		err := rows.Scan(
			&event.ID,
			&event.EventName,
			&event.AllSeats,
			&event.Booked,
			&event.CreatedAt,
			&bookingsJson,
		)
		if err != nil {
			return nil, fmt.Errorf("could not get bookings from db: %w", err)
		}

		if err := json.Unmarshal(bookingsJson, &bookings); err != nil {
			return nil, fmt.Errorf("could not get bookings from db: %w", err)
		}

		event.Bookings = bookings

		events = append(events, event)
	}

	return events, nil
}
