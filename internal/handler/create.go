package handler

import (
	"net/http"
	"strconv"

	"github.com/Komilov31/event-booker/internal/dto"
	_ "github.com/Komilov31/event-booker/internal/model"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// CreateBooking creates a new booking for an event
// @Summary Create a booking
// @Description Create a new booking for a specific event
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param booking body dto.CreateBooking true "Booking data"
// @Success 200 {object} model.Booking
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events/{id}/book [post]
func (h *Handler) CreateBooking(c *ginext.Context) {
	eventID := c.Param("id")
	id, err := strconv.Atoi(eventID)
	if err != nil {
		zlog.Logger.Error().Msg("invalid id" + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	var booking dto.CreateBooking
	booking.EventID = id
	if err := c.BindJSON(&booking); err != nil {
		zlog.Logger.Error().Msg("invalid payload" + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid payload: " + err.Error()})
		return
	}

	book, err := h.service.CreateBooking(booking)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled POST request for creating booking")
	c.JSON(http.StatusOK, book)
}

// CreateEvent creates a new event
// @Summary Create an event
// @Description Create a new event with specified details
// @Tags events
// @Accept json
// @Produce json
// @Param event body dto.CreateEvent true "Event data"
// @Success 200 {object} model.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events [post]
func (h *Handler) CreateEvent(c *ginext.Context) {
	var createEvent dto.CreateEvent
	if err := c.BindJSON(&createEvent); err != nil {
		zlog.Logger.Error().Msg("invalid payload" + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	event, err := h.service.CreateEvent(createEvent)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled POST request for creating event")
	c.JSON(http.StatusOK, event)
}

// ConfirmPayment confirms payment for a booking
// @Summary Confirm payment
// @Description Confirm payment for a specific booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events/{id}/confirm [post]
func (h *Handler) ConfirmPayment(c *ginext.Context) {
	bookingID := c.Param("id")
	id, err := strconv.Atoi(bookingID)
	if err != nil {
		zlog.Logger.Error().Msg("invalid id" + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	if err := h.service.UpdateBookingStatus(id, "paid"); err != nil {
		zlog.Logger.Error().Msg("could not confirm payment: " + err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "could not confirm payment: " + err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled POST request for confirming payment")
	c.JSON(http.StatusOK, ginext.H{"status": "payment confirmed"})
}
