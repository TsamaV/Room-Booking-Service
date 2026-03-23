package schedule_test

import (
	"Avito/internal/schedule"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func NewScheduleHandler2(service *schedule.ScheduleService) *schedule.ScheduleHandler {
	return &schedule.ScheduleHandler{ScheduleService: service}
}

type MockBookingRepo struct{}

func (m *MockBookingRepo) IsSlotBooked(slotID uuid.UUID) (bool, error) {
	return false, nil
}

func TestScheduleHandlerCreate(t *testing.T) {
	scheduleRepo := &MockScheduleRepo{}
	slotRepo := &MockSlotRepo{}
	service := schedule.NewScheduleService(scheduleRepo, slotRepo, &MockBookingRepo{})
	handler := schedule.NewScheduleHandler2(service)

	body, _ := json.Marshal(map[string]any{
		"day_of_week": []int{1, 2, 3},
		"start_time":  "09:00",
		"end_time":    "18:00",
	})

	req := httptest.NewRequest("POST", "/rooms/test/schedule/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("roomId", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestScheduleHandlerGetSlots(t *testing.T) {
	scheduleRepo := &MockScheduleRepo{}
	slotRepo := &MockSlotRepo{}
	service := schedule.NewScheduleService(scheduleRepo, slotRepo, &MockBookingRepo{})
	handler := schedule.NewScheduleHandler2(service)

	req := httptest.NewRequest("GET", "/rooms/test/slots/list?date=2026-03-23", nil)
	req.SetPathValue("roomId", "00000000-0000-0000-0000-000000000001")
	w := httptest.NewRecorder()

	handler.GetSlots(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
