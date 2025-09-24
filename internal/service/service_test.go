package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Komilov31/event-booker/internal/dto"
	"github.com/Komilov31/event-booker/internal/model"
)

type mockStorage struct {
	getEventByIDFunc             func(id int) (*model.Event, error)
	getBookingByIDFunc           func(id int) (*dto.BookingDTO, error)
	getEventWithBookingsByIDFunc func(id int) (*model.Event, error)
	getAllEventsFunc             func() ([]model.Event, error)
	createBookingFunc            func(booking dto.CreateBooking) (*model.Booking, error)
	createEventFunc              func(event dto.CreateEvent) (*model.Event, error)
	updateBookingStatusFunc      func(id int, newStatus string) error
	deleteEventFunc              func(id int) error
	deleteBookingFunc            func(id int, placeCount int) error
}

func (m *mockStorage) GetEventByID(id int) (*model.Event, error) {
	if m.getEventByIDFunc != nil {
		return m.getEventByIDFunc(id)
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockStorage) GetBookingByID(id int) (*dto.BookingDTO, error) {
	if m.getBookingByIDFunc != nil {
		return m.getBookingByIDFunc(id)
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockStorage) GetEventWithBookingsByID(id int) (*model.Event, error) {
	if m.getEventWithBookingsByIDFunc != nil {
		return m.getEventWithBookingsByIDFunc(id)
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockStorage) GetAllEvents() ([]model.Event, error) {
	if m.getAllEventsFunc != nil {
		return m.getAllEventsFunc()
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockStorage) CreateBooking(booking dto.CreateBooking) (*model.Booking, error) {
	if m.createBookingFunc != nil {
		return m.createBookingFunc(booking)
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockStorage) CreateEvent(event dto.CreateEvent) (*model.Event, error) {
	if m.createEventFunc != nil {
		return m.createEventFunc(event)
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockStorage) UpdateBookingStatus(id int, newStatus string) error {
	if m.updateBookingStatusFunc != nil {
		return m.updateBookingStatusFunc(id, newStatus)
	}
	return errors.New("mock not implemented")
}

func (m *mockStorage) DeleteEvent(id int) error {
	if m.deleteEventFunc != nil {
		return m.deleteEventFunc(id)
	}
	return errors.New("mock not implemented")
}

func (m *mockStorage) DeleteBooking(id int, placeCount int) error {
	if m.deleteBookingFunc != nil {
		return m.deleteBookingFunc(id, placeCount)
	}
	return errors.New("mock not implemented")
}

type mockQueue struct {
	publishFunc func(booking dto.QueueMessage) error
	consumeFunc func(ctx context.Context) (<-chan []byte, error)
}

func (m *mockQueue) Publish(booking dto.QueueMessage) error {
	if m.publishFunc != nil {
		return m.publishFunc(booking)
	}
	return errors.New("mock not implemented")
}

func (m *mockQueue) Consume(ctx context.Context) (<-chan []byte, error) {
	if m.consumeFunc != nil {
		return m.consumeFunc(ctx)
	}
	ch := make(chan []byte)
	close(ch)
	return ch, errors.New("mock not implemented")
}

type mockSender struct {
	sendToTelegramFunc func(telegramId int, text string) error
}

func (m *mockSender) SendToTelegram(telegramId int, text string) error {
	if m.sendToTelegramFunc != nil {
		return m.sendToTelegramFunc(telegramId, text)
	}
	return errors.New("mock not implemented")
}

func TestService_CreateBooking_Success(t *testing.T) {
	mockStorage := &mockStorage{
		createBookingFunc: func(booking dto.CreateBooking) (*model.Booking, error) {
			return &model.Booking{
				ID:          1,
				EventID:     booking.EventID,
				PlacesCount: booking.PlacesCount,
				Status:      "pending",
				TelegramID:  booking.TelegramID,
				CreatedAt:   time.Now(),
			}, nil
		},
	}

	mockQueue := &mockQueue{
		publishFunc: func(booking dto.QueueMessage) error {
			return nil
		},
	}

	s := New(mockStorage, mockQueue, &mockSender{})

	booking := dto.CreateBooking{
		EventID:     1,
		TelegramID:  123,
		PlacesCount: 2,
	}

	result, err := s.CreateBooking(booking)

	if err != nil {
		t.Errorf("CreateBooking failed: %v", err)
	}

	if result == nil {
		t.Error("Expected booking result, got nil")
	}

	if result.ID != 1 {
		t.Errorf("Expected booking ID 1, got %d", result.ID)
	}
}

func TestService_CreateBooking_StorageError(t *testing.T) {
	mockStorage := &mockStorage{
		createBookingFunc: func(booking dto.CreateBooking) (*model.Booking, error) {
			return nil, errors.New("storage error")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	booking := dto.CreateBooking{
		EventID:     1,
		TelegramID:  123,
		PlacesCount: 2,
	}

	result, err := s.CreateBooking(booking)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestService_CreateBooking_QueueError(t *testing.T) {
	mockStorage := &mockStorage{
		createBookingFunc: func(booking dto.CreateBooking) (*model.Booking, error) {
			return &model.Booking{
				ID:          1,
				EventID:     booking.EventID,
				PlacesCount: booking.PlacesCount,
				Status:      "pending",
				TelegramID:  booking.TelegramID,
				CreatedAt:   time.Now(),
			}, nil
		},
	}

	mockQueue := &mockQueue{
		publishFunc: func(booking dto.QueueMessage) error {
			return errors.New("queue error")
		},
	}

	s := New(mockStorage, mockQueue, &mockSender{})

	booking := dto.CreateBooking{
		EventID:     1,
		TelegramID:  123,
		PlacesCount: 2,
	}

	result, err := s.CreateBooking(booking)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestService_CreateEvent_Success(t *testing.T) {
	mockStorage := &mockStorage{
		createEventFunc: func(event dto.CreateEvent) (*model.Event, error) {
			return &model.Event{
				ID:        1,
				EventName: event.EventName,
				AllSeats:  event.AllSeats,
				Booked:    0,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	event := dto.CreateEvent{
		EventName: "Test Event",
		AllSeats:  100,
	}

	result, err := s.CreateEvent(event)

	if err != nil {
		t.Errorf("CreateEvent failed: %v", err)
	}

	if result == nil {
		t.Error("Expected event result, got nil")
	}

	if result.EventName != "Test Event" {
		t.Errorf("Expected event name 'Test Event', got '%s'", result.EventName)
	}
}

func TestService_CreateEvent_Error(t *testing.T) {
	mockStorage := &mockStorage{
		createEventFunc: func(event dto.CreateEvent) (*model.Event, error) {
			return nil, errors.New("storage error")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	event := dto.CreateEvent{
		EventName: "Test Event",
		AllSeats:  100,
	}

	result, err := s.CreateEvent(event)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestService_GetEventByID_Success(t *testing.T) {
	expectedEvent := &model.Event{
		ID:        1,
		EventName: "Test Event",
		AllSeats:  100,
		Booked:    0,
		CreatedAt: time.Now(),
	}

	mockStorage := &mockStorage{
		getEventByIDFunc: func(id int) (*model.Event, error) {
			return expectedEvent, nil
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	result, err := s.GetEventByID(1)

	if err != nil {
		t.Errorf("GetEventByID failed: %v", err)
	}

	if result == nil {
		t.Error("Expected event result, got nil")
	}

	if result.ID != 1 {
		t.Errorf("Expected event ID 1, got %d", result.ID)
	}
}

func TestService_GetEventByID_Error(t *testing.T) {
	mockStorage := &mockStorage{
		getEventByIDFunc: func(id int) (*model.Event, error) {
			return nil, errors.New("event not found")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	result, err := s.GetEventByID(1)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestService_GetBookingByID_Success(t *testing.T) {
	expectedBooking := &dto.BookingDTO{
		ID:         1,
		EventID:    1,
		EventName:  "Test Event",
		TelegramID: 123,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	mockStorage := &mockStorage{
		getBookingByIDFunc: func(id int) (*dto.BookingDTO, error) {
			return expectedBooking, nil
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	result, err := s.GetBookingByID(1)

	if err != nil {
		t.Errorf("GetBookingByID failed: %v", err)
	}

	if result == nil {
		t.Error("Expected booking result, got nil")
	}

	if result.ID != 1 {
		t.Errorf("Expected booking ID 1, got %d", result.ID)
	}
}

func TestService_GetBookingByID_Error(t *testing.T) {
	mockStorage := &mockStorage{
		getBookingByIDFunc: func(id int) (*dto.BookingDTO, error) {
			return nil, errors.New("booking not found")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	result, err := s.GetBookingByID(1)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestService_GetEventWithBookingsByID_Success(t *testing.T) {
	expectedEvent := &model.Event{
		ID:        1,
		EventName: "Test Event",
		AllSeats:  100,
		Booked:    10,
		CreatedAt: time.Now(),
		Bookings: []model.Booking{
			{
				ID:          1,
				EventID:     1,
				PlacesCount: 2,
				Status:      "paid",
				TelegramID:  123,
				CreatedAt:   time.Now(),
			},
		},
	}

	mockStorage := &mockStorage{
		getEventWithBookingsByIDFunc: func(id int) (*model.Event, error) {
			return expectedEvent, nil
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	result, err := s.GetEventWithBookingsByID(1)

	if err != nil {
		t.Errorf("GetEventWithBookingsByID failed: %v", err)
	}

	if result == nil {
		t.Error("Expected event result, got nil")
	}

	if result.ID != 1 {
		t.Errorf("Expected event ID 1, got %d", result.ID)
	}

	if len(result.Bookings) != 1 {
		t.Errorf("Expected 1 booking, got %d", len(result.Bookings))
	}
}

func TestService_GetEventWithBookingsByID_Error(t *testing.T) {
	mockStorage := &mockStorage{
		getEventWithBookingsByIDFunc: func(id int) (*model.Event, error) {
			return nil, errors.New("event not found")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	result, err := s.GetEventWithBookingsByID(1)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestService_GetAllEvents_Success(t *testing.T) {
	expectedEvents := []model.Event{
		{
			ID:        1,
			EventName: "Event 1",
			AllSeats:  100,
			Booked:    10,
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			EventName: "Event 2",
			AllSeats:  50,
			Booked:    5,
			CreatedAt: time.Now(),
		},
	}

	mockStorage := &mockStorage{
		getAllEventsFunc: func() ([]model.Event, error) {
			return expectedEvents, nil
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	result, err := s.GetAllEvents()

	if err != nil {
		t.Errorf("GetAllEvents failed: %v", err)
	}

	if result == nil {
		t.Error("Expected events result, got nil")
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 events, got %d", len(result))
	}
}

func TestService_GetAllEvents_Error(t *testing.T) {
	mockStorage := &mockStorage{
		getAllEventsFunc: func() ([]model.Event, error) {
			return nil, errors.New("database error")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	result, err := s.GetAllEvents()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestService_UpdateBookingStatus_Success(t *testing.T) {
	mockStorage := &mockStorage{
		updateBookingStatusFunc: func(id int, newStatus string) error {
			return nil
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	err := s.UpdateBookingStatus(1, "paid")

	if err != nil {
		t.Errorf("UpdateBookingStatus failed: %v", err)
	}
}

func TestService_UpdateBookingStatus_Error(t *testing.T) {
	mockStorage := &mockStorage{
		updateBookingStatusFunc: func(id int, newStatus string) error {
			return errors.New("booking not found")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	err := s.UpdateBookingStatus(1, "paid")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestService_DeleteEvent_Success(t *testing.T) {
	mockStorage := &mockStorage{
		deleteEventFunc: func(id int) error {
			return nil
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	err := s.DeleteEvent(1)

	if err != nil {
		t.Errorf("DeleteEvent failed: %v", err)
	}
}

func TestService_DeleteEvent_Error(t *testing.T) {
	mockStorage := &mockStorage{
		deleteEventFunc: func(id int) error {
			return errors.New("event not found")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	err := s.DeleteEvent(1)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestService_DeleteBooking_Success(t *testing.T) {
	mockStorage := &mockStorage{
		deleteBookingFunc: func(id int, placeCount int) error {
			return nil
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	err := s.DeleteBooking(1, 2)

	if err != nil {
		t.Errorf("DeleteBooking failed: %v", err)
	}
}

func TestService_DeleteBooking_Error(t *testing.T) {
	mockStorage := &mockStorage{
		deleteBookingFunc: func(id int, placeCount int) error {
			return errors.New("booking not found")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	err := s.DeleteBooking(1, 2)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestService_handleQueueMessage_Success(t *testing.T) {
	bookingInfo := &dto.BookingDTO{
		ID:         1,
		EventID:    1,
		EventName:  "Test Event",
		TelegramID: 123,
		Status:     "unpaid",
		CreatedAt:  time.Now(),
	}

	mockStorage := &mockStorage{
		getBookingByIDFunc: func(id int) (*dto.BookingDTO, error) {
			return bookingInfo, nil
		},
		deleteBookingFunc: func(id int, placeCount int) error {
			return nil
		},
	}

	mockSender := &mockSender{
		sendToTelegramFunc: func(telegramId int, text string) error {
			return nil
		},
	}

	s := New(mockStorage, &mockQueue{}, mockSender)

	messageData := `{"booking_id":1,"places_count":2}`

	err := s.handleQueueMessage([]byte(messageData))

	if err != nil {
		t.Errorf("handleQueueMessage failed: %v", err)
	}
}

func TestService_handleQueueMessage_UnmarshalError(t *testing.T) {
	s := New(&mockStorage{}, &mockQueue{}, &mockSender{})

	invalidMessageData := `{"invalid_json"}`

	err := s.handleQueueMessage([]byte(invalidMessageData))

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestService_handleQueueMessage_GetBookingError(t *testing.T) {
	mockStorage := &mockStorage{
		getBookingByIDFunc: func(id int) (*dto.BookingDTO, error) {
			return nil, errors.New("booking not found")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	messageData := `{"booking_id":1,"places_count":2}`

	err := s.handleQueueMessage([]byte(messageData))

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestService_handleQueueMessage_PaidBooking(t *testing.T) {
	bookingInfo := &dto.BookingDTO{
		ID:         1,
		EventID:    1,
		EventName:  "Test Event",
		TelegramID: 123,
		Status:     "paid",
		CreatedAt:  time.Now(),
	}

	mockStorage := &mockStorage{
		getBookingByIDFunc: func(id int) (*dto.BookingDTO, error) {
			return bookingInfo, nil
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	messageData := `{"booking_id":1,"places_count":2}`

	err := s.handleQueueMessage([]byte(messageData))

	if err != nil {
		t.Errorf("handleQueueMessage failed for paid booking: %v", err)
	}
}

func TestService_handleQueueMessage_DeleteBookingError(t *testing.T) {
	bookingInfo := &dto.BookingDTO{
		ID:         1,
		EventID:    1,
		EventName:  "Test Event",
		TelegramID: 123,
		Status:     "unpaid",
		CreatedAt:  time.Now(),
	}

	mockStorage := &mockStorage{
		getBookingByIDFunc: func(id int) (*dto.BookingDTO, error) {
			return bookingInfo, nil
		},
		deleteBookingFunc: func(id int, placeCount int) error {
			return errors.New("delete booking error")
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	messageData := `{"booking_id":1,"places_count":2}`

	err := s.handleQueueMessage([]byte(messageData))

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestService_handleQueueMessage_SendTelegramError(t *testing.T) {
	bookingInfo := &dto.BookingDTO{
		ID:         1,
		EventID:    1,
		EventName:  "Test Event",
		TelegramID: 123,
		Status:     "unpaid",
		CreatedAt:  time.Now(),
	}

	mockStorage := &mockStorage{
		getBookingByIDFunc: func(id int) (*dto.BookingDTO, error) {
			return bookingInfo, nil
		},
		deleteBookingFunc: func(id int, placeCount int) error {
			return nil
		},
	}

	mockSender := &mockSender{
		sendToTelegramFunc: func(telegramId int, text string) error {
			return errors.New("telegram send error")
		},
	}

	s := New(mockStorage, &mockQueue{}, mockSender)

	messageData := `{"booking_id":1,"places_count":2}`

	err := s.handleQueueMessage([]byte(messageData))

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestService_handleQueueMessage_NoTelegramID(t *testing.T) {
	bookingInfo := &dto.BookingDTO{
		ID:         1,
		EventID:    1,
		EventName:  "Test Event",
		TelegramID: 0,
		Status:     "unpaid",
		CreatedAt:  time.Now(),
	}

	mockStorage := &mockStorage{
		getBookingByIDFunc: func(id int) (*dto.BookingDTO, error) {
			return bookingInfo, nil
		},
		deleteBookingFunc: func(id int, placeCount int) error {
			return nil
		},
	}

	s := New(mockStorage, &mockQueue{}, &mockSender{})

	messageData := `{"booking_id":1,"places_count":2}`

	err := s.handleQueueMessage([]byte(messageData))

	if err != nil {
		t.Errorf("handleQueueMessage failed for booking without telegram ID: %v", err)
	}
}
