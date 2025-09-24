package repository

import "fmt"

func (r *Repository) DeleteEvent(id int) error {
	tx, err := r.db.Master.Begin()
	if err != nil {
		return fmt.Errorf("could not start trascation: %w", err)
	}
	defer tx.Rollback()

	query := "DELETE FROM bookings WHERE event_id = $1"
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not delete event: %w", err)
	}

	query = "DELETE FROM events WHERE id = $1"
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not delete event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transcation: %w", err)
	}

	return nil
}

func (r *Repository) DeleteBooking(id int, placeCount int) error {
	tx, err := r.db.Master.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	query := "DELETE FROM bookings WHERE id = $1 RETURNING event_id"

	var eventID int
	err = tx.QueryRow(query, id).Scan(&eventID)
	if err != nil {
		return fmt.Errorf("could not delete booking: %w", err)
	}

	query = `UPDATE events
	SET booked = booked - $1
	WHERE id = $2`
	result, err := tx.Exec(query, placeCount, eventID)
	if err != nil {
		return fmt.Errorf("could decrement bookings in event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not decrement booking in event: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNoSuchEvent
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transcation: %w", err)
	}

	return nil
}
