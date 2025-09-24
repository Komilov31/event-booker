package handler

import (
	"net/http"
	"strconv"

	_ "github.com/Komilov31/event-booker/internal/dto"
	_ "github.com/Komilov31/event-booker/internal/model"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// GetBooking retrieves booking information by ID
// @Summary Get booking by ID
// @Description Get detailed information about a specific booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} dto.BookingDTO
// @Failure 400 {object} map[string]string
// @Router /events/{id} [get]
func (h *Handler) GetBooking(c *ginext.Context) {
	bookingID := c.Param("id")
	id, err := strconv.Atoi(bookingID)
	if err != nil {
		zlog.Logger.Error().Msg("invalid id" + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	booking, err := h.service.GetBookingByID(id)
	if err != nil {
		zlog.Logger.Error().Msg("could not get booking info: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled GET request and retuned booking info")
	c.JSON(http.StatusOK, booking)
}

// GetEvent retrieves event information by ID
// @Summary Get event by ID
// @Description Get detailed information about a specific event
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} model.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events/{id} [get]
func (h *Handler) GetEvent(c *ginext.Context) {
	eventID := c.Param("id")
	id, err := strconv.Atoi(eventID)
	if err != nil {
		zlog.Logger.Error().Msg("invalid id" + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	event, err := h.service.GetEventByID(id)
	if err != nil {
		zlog.Logger.Error().Msg("could not get event info: " + err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled GET request and retuned event info")
	c.JSON(http.StatusOK, event)
}

// GetAllEvents retrieves all events
// @Summary Get all events
// @Description Get a list of all available events
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {array} model.Event
// @Failure 500 {object} map[string]string
// @Router /events [get]
func (h *Handler) GetAllEvents(c *ginext.Context) {
	events, err := h.service.GetAllEvents()
	if err != nil {
		zlog.Logger.Error().Msg("could not get all events: " + err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled GET request and retuned all events")
	c.JSON(http.StatusOK, events)
}

// GetMainPage godoc
// @Summary      Get main page
// @Description  Get the main HTML page of the application
// @Tags         pages
// @Accept       json
// @Produce      html
// @Success      200  {string} string "HTML page content"
// @Router       / [get]
func (h *Handler) GetMainPage(c *ginext.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// GetAdminPage godoc
// @Summary      Get admin page
// @Description  Get the main HTML page of the application
// @Tags         pages
// @Accept       json
// @Produce      html
// @Success      200  {string} string "HTML page content"
// @Router       / [get]
func (h *Handler) GetAdminPage(c *ginext.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

// GetUserPage godoc
// @Summary      Get user page
// @Description  Get the user HTML page of the application
// @Tags         pages
// @Accept       json
// @Produce      html
// @Success      200  {string} string "HTML page content"
// @Router       / [get]
func (h *Handler) GetUserPage(c *ginext.Context) {
	c.HTML(http.StatusOK, "user.html", nil)
}
