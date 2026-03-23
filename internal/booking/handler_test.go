package booking_test

import (
	"Avito/internal/booking"
	"Avito/pkg/middleware"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestHandlerCreate(t *testing.T) {
	userID := uuid.New()
	slotID := uuid.New()

	repo := &MockBookingRepo{}
	slotRepo := &MockSlotRepo{}
	service := booking.NewBookingService(repo, slotRepo)
	handler := booking.NewBookingHandler2(service)

	body, _ := json.Marshal(map[string]string{
		"slot_id": slotID.String(),
	})

	req := httptest.NewRequest("POST", "/bookings/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middleware.ContextUserID, userID.String())
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestHandlerCancel(t *testing.T) {
	userID := uuid.New()
	bookingID := uuid.New()

	repo := &MockBookingRepoWithBooking{
		booking: &booking.Booking{
			ID:     bookingID,
			UserID: userID,
			Status: "active",
		},
	}

	service := booking.NewBookingService(repo, &MockSlotRepo{})
	handler := booking.NewBookingHandler2(service)

	req := httptest.NewRequest("POST", "/bookings/"+bookingID.String()+"/cancel", nil)
	req.SetPathValue("bookingId", bookingID.String())

	ctx := context.WithValue(req.Context(), middleware.ContextUserID, userID.String())
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.Cancel(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandlerGetAll(t *testing.T) {
	userID := uuid.New()

	repo := &MockBookingRepoWithBooking{
		booking: &booking.Booking{
			ID:     uuid.New(),
			UserID: userID,
			Status: "active",
		},
	}

	service := booking.NewBookingService(repo, &MockSlotRepo{})
	handler := booking.NewBookingHandler2(service)

	req := httptest.NewRequest("GET", "/bookings/list", nil)
	w := httptest.NewRecorder()
	handler.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandlerGetMy(t *testing.T) {
	userID := uuid.New()

	repo := &MockBookingRepoWithBooking{
		booking: &booking.Booking{
			ID:     uuid.New(),
			UserID: userID,
			Status: "active",
		},
	}

	service := booking.NewBookingService(repo, &MockSlotRepo{})
	handler := booking.NewBookingHandler2(service)

	req := httptest.NewRequest("GET", "/bookings/my", nil)
	ctx := context.WithValue(req.Context(), middleware.ContextUserID, userID.String())
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetMy(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}