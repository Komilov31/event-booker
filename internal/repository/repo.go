package repository

import (
	"errors"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrNoSuchEvent      = errors.New("there is not event with such id")
	ErrNoSuchBooking    = errors.New("there is not booking with such id")
	ErrNoSeatsAvailable = errors.New("all seats are booked")
)

type Repository struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Repository {
	return &Repository{
		db: db,
	}
}
