package booking

import (
	"Avito/pkg/middleware"
	"Avito/pkg/req"
	"Avito/pkg/res"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type BookingHandler struct {
	BookingService *BookingService
}

func NewBookingHandler(router *http.ServeMux, service *BookingService, jwtSecret string) {
	handler := &BookingHandler{
		BookingService: service,
	}

	router.Handle("POST /bookings/create", middleware.IsAuthed(middleware.UserOnly(http.HandlerFunc(handler.Create)), jwtSecret))
	router.Handle("GET /bookings/my", middleware.IsAuthed(middleware.UserOnly(http.HandlerFunc(handler.GetMy)), jwtSecret))
	router.Handle("POST /bookings/{bookingId}/cancel", middleware.IsAuthed(middleware.UserOnly(http.HandlerFunc(handler.Cancel)), jwtSecret))
	router.Handle("GET /bookings/list", middleware.IsAuthed(middleware.AdminOnly(http.HandlerFunc(handler.GetAll)), jwtSecret))
}

func NewBookingHandler2(service *BookingService) *BookingHandler {
	return &BookingHandler{BookingService: service}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.ContextUserID).(string)
	if !ok {
		res.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		res.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user id")
		return
	}

	body, err := req.Decode[CreateBookingRequest](r)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	slotID, err := uuid.Parse(body.SlotID)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid slot id")
		return
	}

	booking, err := h.BookingService.Create(userID, slotID)
	if err != nil {
		if err.Error() == "slot is already booked" {
			res.Error(w, http.StatusConflict, "SLOT_ALREADY_BOOKED", err.Error())
			return
		}
		if err.Error() == "slot not found" {
			res.Error(w, http.StatusNotFound, "SLOT_NOT_FOUND", err.Error())
			return
		}
		if err.Error() == "cannot book a slot in the past" {
        res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    	}
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	res.JSON(w, http.StatusCreated, map[string]any{
		"booking": BookingResponse{
			ID:             booking.ID.String(),
			SlotID:         booking.SlotID.String(),
			UserID:         booking.UserID.String(),
			Status:         booking.Status,
			ConferenceLink: booking.ConferenceLink,
			CreatedAt:      booking.CreatedAt.String(),
		},
	})
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.ContextUserID).(string)
	if !ok {
		res.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		res.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user id")
		return
	}

	idStr := r.PathValue("bookingId")
	bookingID, err := uuid.Parse(idStr)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid booking id")
		return
	}

	booking, err := h.BookingService.Cancel(bookingID, userID)
	if err != nil {
		if err.Error() == "booking not found" {
			res.Error(w, http.StatusNotFound, "BOOKING_NOT_FOUND", err.Error())
			return
		}
		if err.Error() == "you can only cancel your own bookings" {
			res.Error(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		res.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	res.JSON(w, http.StatusOK, map[string]any{
		"booking": BookingResponse{
			ID:             booking.ID.String(),
			SlotID:         booking.SlotID.String(),
			UserID:         booking.UserID.String(),
			Status:         booking.Status,
			ConferenceLink: booking.ConferenceLink,
			CreatedAt:      booking.CreatedAt.String(),
		},
	})
}

func (h *BookingHandler) GetMy(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.ContextUserID).(string)
	if !ok {
		res.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		res.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user id")
		return
	}

	bookings, err := h.BookingService.GetMyBookings(userID)
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	var response []BookingResponse
	for _, b := range bookings {
		response = append(response, BookingResponse{
			ID:             b.ID.String(),
			SlotID:         b.SlotID.String(),
			UserID:         b.UserID.String(),
			Status:         b.Status,
			ConferenceLink: b.ConferenceLink,
			CreatedAt:      b.CreatedAt.String(),
		})
	}

	res.JSON(w, http.StatusOK, map[string]any{
		"bookings": response,
	})
}

func (h *BookingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	bookings, total, err := h.BookingService.GetAllPaginated(page, pageSize)
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	var data []BookingResponse
	for _, b := range bookings {
		data = append(data, BookingResponse{
			ID:             b.ID.String(),
			SlotID:         b.SlotID.String(),
			UserID:         b.UserID.String(),
			Status:         b.Status,
			ConferenceLink: b.ConferenceLink,
			CreatedAt:      b.CreatedAt.String(),
		})
	}

	res.JSON(w, http.StatusOK, map[string]any{
		"bookings": data,
		"pagination": map[string]any{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	})
}