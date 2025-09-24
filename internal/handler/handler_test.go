package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Komilov31/event-booker/internal/dto"
	"github.com/Komilov31/event-booker/internal/model"
	"github.com/gin-gonic/gin"
)

type mockEventBookerService struct {
	getEventByIDFunc        func(id int) (*model.Event, error)
	getBookingByIDFunc      func(id int) (*dto.BookingDTO, error)
	getAllEventsFunc        func() ([]model.Event, error)
	createBookingFunc       func(booking dto.CreateBooking) (*model.Booking, error)
	createEventFunc         func(event dto.CreateEvent) (*model.Event, error)
	updateBookingStatusFunc func(id int, newStatus string) error
}

func (m *mockEventBookerService) GetEventByID(id int) (*model.Event, error) {
	if m.getEventByIDFunc != nil {
		return m.getEventByIDFunc(id)
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockEventBookerService) GetBookingByID(id int) (*dto.BookingDTO, error) {
	if m.getBookingByIDFunc != nil {
		return m.getBookingByIDFunc(id)
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockEventBookerService) GetAllEvents() ([]model.Event, error) {
	if m.getAllEventsFunc != nil {
		return m.getAllEventsFunc()
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockEventBookerService) CreateBooking(booking dto.CreateBooking) (*model.Booking, error) {
	if m.createBookingFunc != nil {
		return m.createBookingFunc(booking)
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockEventBookerService) CreateEvent(event dto.CreateEvent) (*model.Event, error) {
	if m.createEventFunc != nil {
		return m.createEventFunc(event)
	}
	return nil, errors.New("mock not implemented")
}

func (m *mockEventBookerService) UpdateBookingStatus(id int, newStatus string) error {
	if m.updateBookingStatusFunc != nil {
		return m.updateBookingStatusFunc(id, newStatus)
	}
	return errors.New("mock not implemented")
}

func createMockContext(method, path string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	var jsonBody []byte
	if body != nil {
		jsonBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c, w
}

func TestHandler_CreateBooking_Success(t *testing.T) {
	expectedBooking := &model.Booking{
		ID:          1,
		EventID:     1,
		PlacesCount: 2,
		Status:      "pending",
		TelegramID:  123,
	}

	mockService := &mockEventBookerService{
		createBookingFunc: func(booking dto.CreateBooking) (*model.Booking, error) {
			return expectedBooking, nil
		},
	}

	h := New(mockService)

	booking := dto.CreateBooking{
		EventID:     1,
		TelegramID:  123,
		PlacesCount: 2,
	}

	ctx, w := createMockContext("POST", "/events/1/book", booking)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	h.CreateBooking(ctx)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response model.Booking
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.ID != 1 {
		t.Errorf("Expected booking ID 1, got %d", response.ID)
	}
}

func TestHandler_CreateBooking_InvalidID(t *testing.T) {
	mockService := &mockEventBookerService{}

	h := New(mockService)

	booking := dto.CreateBooking{
		EventID:     1,
		TelegramID:  123,
		PlacesCount: 2,
	}

	ctx, w := createMockContext("POST", "/events/abc/book", booking)
	ctx.Params = []gin.Param{{Key: "id", Value: "abc"}}

	h.CreateBooking(ctx)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_CreateBooking_InvalidPayload(t *testing.T) {
	mockService := &mockEventBookerService{}

	h := New(mockService)

	invalidBody := `{"invalid": json}`

	ctx, w := createMockContext("POST", "/events/1/book", nil)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}
	ctx.Request.Body = &mockBody{data: invalidBody}

	h.CreateBooking(ctx)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_CreateBooking_ServiceError(t *testing.T) {
	mockService := &mockEventBookerService{
		createBookingFunc: func(booking dto.CreateBooking) (*model.Booking, error) {
			return nil, errors.New("service error")
		},
	}

	h := New(mockService)

	booking := dto.CreateBooking{
		EventID:     1,
		TelegramID:  123,
		PlacesCount: 2,
	}

	ctx, w := createMockContext("POST", "/events/1/book", booking)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	h.CreateBooking(ctx)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestHandler_CreateEvent_Success(t *testing.T) {
	expectedEvent := &model.Event{
		ID:        1,
		EventName: "Test Event",
		AllSeats:  100,
		Booked:    0,
	}

	mockService := &mockEventBookerService{
		createEventFunc: func(event dto.CreateEvent) (*model.Event, error) {
			return expectedEvent, nil
		},
	}

	h := New(mockService)

	event := dto.CreateEvent{
		EventName: "Test Event",
		AllSeats:  100,
	}

	ctx, w := createMockContext("POST", "/events", event)

	h.CreateEvent(ctx)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response model.Event
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.EventName != "Test Event" {
		t.Errorf("Expected event name 'Test Event', got '%s'", response.EventName)
	}
}

func TestHandler_CreateEvent_InvalidPayload(t *testing.T) {
	mockService := &mockEventBookerService{}

	h := New(mockService)

	invalidBody := `{"invalid": json}`

	ctx, w := createMockContext("POST", "/events", nil)
	ctx.Request.Body = &mockBody{data: invalidBody}

	h.CreateEvent(ctx)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_CreateEvent_ServiceError(t *testing.T) {
	mockService := &mockEventBookerService{
		createEventFunc: func(event dto.CreateEvent) (*model.Event, error) {
			return nil, errors.New("service error")
		},
	}

	h := New(mockService)

	event := dto.CreateEvent{
		EventName: "Test Event",
		AllSeats:  100,
	}

	ctx, w := createMockContext("POST", "/events", event)

	h.CreateEvent(ctx)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestHandler_ConfirmPayment_Success(t *testing.T) {
	mockService := &mockEventBookerService{
		updateBookingStatusFunc: func(id int, newStatus string) error {
			return nil
		},
	}

	h := New(mockService)

	ctx, w := createMockContext("POST", "/events/1/confirm", nil)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	h.ConfirmPayment(ctx)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "payment confirmed" {
		t.Errorf("Expected status 'payment confirmed', got '%s'", response["status"])
	}
}

func TestHandler_ConfirmPayment_InvalidID(t *testing.T) {
	mockService := &mockEventBookerService{}

	h := New(mockService)

	ctx, w := createMockContext("POST", "/events/abc/confirm", nil)
	ctx.Params = []gin.Param{{Key: "id", Value: "abc"}}

	h.ConfirmPayment(ctx)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_ConfirmPayment_ServiceError(t *testing.T) {
	mockService := &mockEventBookerService{
		updateBookingStatusFunc: func(id int, newStatus string) error {
			return errors.New("service error")
		},
	}

	h := New(mockService)

	ctx, w := createMockContext("POST", "/events/1/confirm", nil)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	h.ConfirmPayment(ctx)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestHandler_GetBooking_Success(t *testing.T) {
	expectedBooking := &dto.BookingDTO{
		ID:         1,
		EventID:    1,
		EventName:  "Test Event",
		TelegramID: 123,
		Status:     "pending",
	}

	mockService := &mockEventBookerService{
		getBookingByIDFunc: func(id int) (*dto.BookingDTO, error) {
			return expectedBooking, nil
		},
	}

	h := New(mockService)

	ctx, w := createMockContext("GET", "/events/1", nil)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	h.GetBooking(ctx)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response dto.BookingDTO
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.ID != 1 {
		t.Errorf("Expected booking ID 1, got %d", response.ID)
	}
}

func TestHandler_GetBooking_InvalidID(t *testing.T) {
	mockService := &mockEventBookerService{}

	h := New(mockService)

	ctx, w := createMockContext("GET", "/events/abc", nil)
	ctx.Params = []gin.Param{{Key: "id", Value: "abc"}}

	h.GetBooking(ctx)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_GetBooking_ServiceError(t *testing.T) {
	mockService := &mockEventBookerService{
		getBookingByIDFunc: func(id int) (*dto.BookingDTO, error) {
			return nil, errors.New("service error")
		},
	}

	h := New(mockService)

	ctx, w := createMockContext("GET", "/events/1", nil)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	h.GetBooking(ctx)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_GetEvent_Success(t *testing.T) {
	expectedEvent := &model.Event{
		ID:        1,
		EventName: "Test Event",
		AllSeats:  100,
		Booked:    0,
	}

	mockService := &mockEventBookerService{
		getEventByIDFunc: func(id int) (*model.Event, error) {
			return expectedEvent, nil
		},
	}

	h := New(mockService)

	ctx, w := createMockContext("GET", "/events/1", nil)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	h.GetEvent(ctx)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response model.Event
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.ID != 1 {
		t.Errorf("Expected event ID 1, got %d", response.ID)
	}
}

func TestHandler_GetEvent_InvalidID(t *testing.T) {
	mockService := &mockEventBookerService{}

	h := New(mockService)

	ctx, w := createMockContext("GET", "/events/abc", nil)
	ctx.Params = []gin.Param{{Key: "id", Value: "abc"}}

	h.GetEvent(ctx)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandler_GetEvent_ServiceError(t *testing.T) {
	mockService := &mockEventBookerService{
		getEventByIDFunc: func(id int) (*model.Event, error) {
			return nil, errors.New("service error")
		},
	}

	h := New(mockService)

	ctx, w := createMockContext("GET", "/events/1", nil)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	h.GetEvent(ctx)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestHandler_GetAllEvents_Success(t *testing.T) {
	expectedEvents := []model.Event{
		{
			ID:        1,
			EventName: "Event 1",
			AllSeats:  100,
			Booked:    10,
		},
		{
			ID:        2,
			EventName: "Event 2",
			AllSeats:  50,
			Booked:    5,
		},
	}

	mockService := &mockEventBookerService{
		getAllEventsFunc: func() ([]model.Event, error) {
			return expectedEvents, nil
		},
	}

	h := New(mockService)

	ctx, w := createMockContext("GET", "/events", nil)

	h.GetAllEvents(ctx)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []model.Event
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 events, got %d", len(response))
	}
}

func TestHandler_GetAllEvents_ServiceError(t *testing.T) {
	mockService := &mockEventBookerService{
		getAllEventsFunc: func() ([]model.Event, error) {
			return nil, errors.New("service error")
		},
	}

	h := New(mockService)

	ctx, w := createMockContext("GET", "/events", nil)

	h.GetAllEvents(ctx)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

type mockBody struct {
	data string
}

func (m *mockBody) Read(p []byte) (n int, err error) {
	copy(p, []byte(m.data))
	return len(m.data), nil
}

func (m *mockBody) Close() error {
	return nil
}
