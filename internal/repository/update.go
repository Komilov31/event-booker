package repository

import "fmt"

func (r *Repository) UpdateBookingStatus(id int, newStatus string) error {
	query := `UPDATE bookings
	SET status = $1
	WHERE id = $2`

	result, err := r.db.Master.Exec(query, newStatus, id)
	if err != nil {
		return fmt.Errorf("could not update bookings status: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNoSuchBooking
	}

	return nil
}
