package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Komilov31/event-booker/internal/dto"
	"github.com/wb-go/wbf/zlog"
)

const (
	statusPaid = "paid"
)

func (s *Service) StartWorker(ctx context.Context) {
	zlog.Logger.Info().Msg("successfully started worker to cancel unpaid bookings")
	go func() {
		s.consumeMessages(ctx)
	}()
}

func (s *Service) consumeMessages(ctx context.Context) {
	messages, err := s.queue.Consume(ctx)
	if err != nil {
		zlog.Logger.Error().Msg("could not consume message from queue: " + err.Error())
		return
	}

	for msg := range messages {
		if err := s.handleQueueMessage(msg); err != nil {
			zlog.Logger.Error().Msg("error while handling queue message: " + err.Error())
			continue
		}
		zlog.Logger.Error().Msg("successfully handled message from queue")
	}
}

func (s *Service) handleQueueMessage(msgData []byte) error {
	var message dto.QueueMessage
	if err := json.Unmarshal(msgData, &message); err != nil {
		return fmt.Errorf("could not unmarshal queue message to model: %w", err)
	}

	bookingInfo, err := s.storage.GetBookingByID(message.BookingID)
	if err != nil {
		return err
	}

	if bookingInfo.Status == statusPaid {
		zlog.Logger.Info().Msgf("user paid booking with id: %d. No need to send notif", message.BookingID)
		return nil
	}

	if err := s.storage.DeleteBooking(message.BookingID, message.PlacesCount); err != nil {
		return err
	}

	if bookingInfo.TelegramID != 0 {
		tgMessage := fmt.Sprintf("Your booking to event %s was cancelled due to unpaid staus", bookingInfo.EventName)
		if err := s.sender.SendToTelegram(bookingInfo.TelegramID, tgMessage); err != nil {
			return err
		}
	}

	return nil
}
