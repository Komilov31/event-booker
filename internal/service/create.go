package service

import (
	"github.com/Komilov31/event-booker/internal/dto"
	"github.com/Komilov31/event-booker/internal/model"
	"github.com/wb-go/wbf/zlog"
)

func (s *Service) CreateBooking(booking dto.CreateBooking) (*model.Booking, error) {
	createBooking, err := s.storage.CreateBooking(booking)
	if err != nil {
		return nil, err
	}

	var message dto.QueueMessage
	message.BookingID = createBooking.ID
	message.PlacesCount = createBooking.PlacesCount

	if err := s.queue.Publish(message); err != nil {
		return nil, err
	}
	zlog.Logger.Error().Msg("successfully published message to queue")

	return createBooking, nil
}

func (s *Service) CreateEvent(event dto.CreateEvent) (*model.Event, error) {
	return s.storage.CreateEvent(event)
}
