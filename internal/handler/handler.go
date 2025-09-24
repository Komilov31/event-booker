package handler

import (
	"github.com/Komilov31/event-booker/internal/dto"
	"github.com/Komilov31/event-booker/internal/model"
)

type EventBookerService interface {
	GetEventByID(id int) (*model.Event, error)
	GetBookingByID(id int) (*dto.BookingDTO, error)
	GetAllEvents() ([]model.Event, error)
	CreateBooking(booking dto.CreateBooking) (*model.Booking, error)
	CreateEvent(event dto.CreateEvent) (*model.Event, error)
	UpdateBookingStatus(id int, newStatus string) error
}

type Handler struct {
	service EventBookerService
}

func New(service EventBookerService) *Handler {
	return &Handler{
		service: service,
	}
}
