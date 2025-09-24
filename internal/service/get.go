package service

import (
	"github.com/Komilov31/event-booker/internal/dto"
	"github.com/Komilov31/event-booker/internal/model"
)

func (s *Service) GetEventByID(id int) (*model.Event, error) {
	return s.storage.GetEventByID(id)
}

func (s *Service) GetBookingByID(id int) (*dto.BookingDTO, error) {
	return s.storage.GetBookingByID(id)
}

func (s *Service) GetEventWithBookingsByID(id int) (*model.Event, error) {
	return s.storage.GetEventWithBookingsByID(id)
}

func (s *Service) GetAllEvents() ([]model.Event, error) {
	return s.storage.GetAllEvents()
}
