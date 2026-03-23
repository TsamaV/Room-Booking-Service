package schedule

import (
	"Avito/pkg/middleware"
	"Avito/pkg/req"
	"Avito/pkg/res"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ScheduleHandler struct {
	ScheduleService *ScheduleService
}

func NewScheduleHandler(router *http.ServeMux, service *ScheduleService, jwtSecret string) {
	handler := &ScheduleHandler{
		ScheduleService: service,
	}

	router.Handle("POST /rooms/{roomId}/schedule/create", middleware.IsAuthed(middleware.AdminOnly(http.HandlerFunc(handler.Create)), jwtSecret))
	router.Handle("GET /rooms/{roomId}/slots/list", middleware.IsAuthed(http.HandlerFunc(handler.GetSlots), jwtSecret))
}

func NewScheduleHandler2(service *ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{ScheduleService: service}
}

func (handler *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("roomId")
	roomID, err := uuid.Parse(idStr)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid room id")
		return
	}

	body, err := req.Decode[CreateScheduleRequest](r)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	for _, d := range body.DayOfWeek {
		if d < 1 || d > 7 {
			res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "daysOfWeek values must be between 1 and 7")
			return
		}
	}

	startTime, err := time.Parse("15:04", body.StartTime)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid start_time format, use HH:MM")
		return
	}

	endTime, err := time.Parse("15:04", body.EndTime)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid end_time format, use HH:MM")
		return
	}

	schedule, err := handler.ScheduleService.Create(roomID, body.DayOfWeek, startTime, endTime)
	if err != nil {
		if err.Error() == "schedule already exists for this room" {
			res.Error(w, http.StatusConflict, "SCHEDULE_EXISTS", err.Error())
			return
		}
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	res.JSON(w, http.StatusCreated, map[string]any{
		"schedule": ScheduleResponse{
			ID:        schedule.ID.String(),
			RoomID:    schedule.RoomID.String(),
			DayOfWeek: []int64(schedule.DayOfWeek),
			StartTime: schedule.StartTime.Format("15:04"),
			EndTime:   schedule.EndTime.Format("15:04"),
		},
	})
}

func (handler *ScheduleHandler) GetSlots(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("roomId")
	roomID, err := uuid.Parse(idStr)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid room id")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "date query param is required")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid date format, use YYYY-MM-DD")
		return
	}

	slots, err := handler.ScheduleService.GetSlots(roomID, date)
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	var response []map[string]any
	for _, slot := range slots {
		response = append(response, map[string]any{
			"id":     slot.ID.String(),
			"roomId": slot.RoomID.String(),
			"start":  slot.StartTime.UTC().Format(time.RFC3339),
			"end":    slot.EndTime.UTC().Format(time.RFC3339),
		})
	}

	res.JSON(w, http.StatusOK, map[string]any{
		"slots": response,
	})
}