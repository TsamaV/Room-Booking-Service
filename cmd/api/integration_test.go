package main_test

import (
	"Avito/configs"
	"Avito/internal/auth"
	"Avito/internal/booking"
	"Avito/internal/room"
	"Avito/internal/schedule"
	"Avito/internal/slot"
	"Avito/internal/user"
	"Avito/internal/db"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupApp() http.Handler {
	os.Setenv("DSN", "postgres://postgres:my_pass@localhost:5432/booking?sslmode=disable")
	os.Setenv("SECRET", "supersecretkey")

	conf := configs.LoadConfig()
	database := db.NewDb(conf)

	database.AutoMigrate(
		&user.User{},
		&room.Room{},
		&schedule.Schedule{},
		&slot.Slot{},
		&booking.Booking{},
	)

	router := http.NewServeMux()

	userRepository := user.NewUserRepository(database)
	roomRepository := room.NewRoomRepository(database)
	scheduleRepository := schedule.NewScheduleRepository(database)
	slotRepository := slot.NewSlotRepository(database)
	bookingRepository := booking.NewBookingRepository(database)

	authService := auth.NewAuthService(userRepository, conf.Auth.Secret)
	roomService := room.NewRoomService(roomRepository)
	scheduleService := schedule.NewScheduleService(scheduleRepository, slotRepository, bookingRepository)
	bookingService := booking.NewBookingService(bookingRepository, slotRepository)

	auth.NewAuthHandler(router, authService)
	room.NewRoomHandler(router, roomService, conf.Auth.Secret)
	schedule.NewScheduleHandler(router, scheduleService, conf.Auth.Secret)
	booking.NewBookingHandler(router, bookingService, conf.Auth.Secret)

	return router
}

func getToken(t *testing.T, app http.Handler, role string) string {
	body, _ := json.Marshal(map[string]string{"role": role})
	req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	return response["token"]
}

func TestCreateRoomScheduleBooking(t *testing.T) {
	app := setupApp()

	adminToken := getToken(t, app, "admin")
	userToken := getToken(t, app, "user")

	roomBody, _ := json.Marshal(map[string]any{
		"name":        "Integration Room",
		"description": "Test",
		"capacity":    10,
	})
	req := httptest.NewRequest("POST", "/rooms/create", bytes.NewBuffer(roomBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 for room creation, got %d: %s", w.Code, w.Body.String())
	}

	var roomResp map[string]map[string]string
	json.NewDecoder(w.Body).Decode(&roomResp)
	roomID := roomResp["room"]["id"]

	scheduleBody, _ := json.Marshal(map[string]any{
		"day_of_week": []int{1, 2, 3, 4, 5, 6, 7},
		"start_time":  "09:00",
		"end_time":    "18:00",
	})
	req = httptest.NewRequest("POST", "/rooms/"+roomID+"/schedule/create", bytes.NewBuffer(scheduleBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.SetPathValue("roomId", roomID)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 for schedule creation, got %d: %s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest("GET", "/rooms/"+roomID+"/slots/list?date=2026-03-29", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.SetPathValue("roomId", roomID)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	var slotsResp map[string][]map[string]string
	json.NewDecoder(w.Body).Decode(&slotsResp)

	if len(slotsResp["slots"]) == 0 {
		t.Fatal("expected slots, got empty")
	}

	slotID := slotsResp["slots"][0]["id"]

	bookingBody, _ := json.Marshal(map[string]string{
		"slot_id": slotID,
	})
	req = httptest.NewRequest("POST", "/bookings/create", bytes.NewBuffer(bookingBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 for booking creation, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCancelBooking(t *testing.T) {
	app := setupApp()

	adminToken := getToken(t, app, "admin")
	userToken := getToken(t, app, "user")

	roomBody, _ := json.Marshal(map[string]any{
		"name":     "Cancel Room",
		"capacity": 5,
	})
	req := httptest.NewRequest("POST", "/rooms/create", bytes.NewBuffer(roomBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	var roomResp map[string]map[string]string
	json.NewDecoder(w.Body).Decode(&roomResp)
	roomID := roomResp["room"]["id"]

	scheduleBody, _ := json.Marshal(map[string]any{
		"day_of_week": []int{1, 2, 3, 4, 5, 6, 7},
		"start_time":  "09:00",
		"end_time":    "18:00",
	})
	req = httptest.NewRequest("POST", "/rooms/"+roomID+"/schedule/create", bytes.NewBuffer(scheduleBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.SetPathValue("roomId", roomID)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	req = httptest.NewRequest("GET", "/rooms/"+roomID+"/slots/list?date=2026-03-29", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.SetPathValue("roomId", roomID)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	var slotsResp map[string][]map[string]string
	json.NewDecoder(w.Body).Decode(&slotsResp)
	slotID := slotsResp["slots"][0]["id"]

	bookingBody, _ := json.Marshal(map[string]string{"slot_id": slotID})
	req = httptest.NewRequest("POST", "/bookings/create", bytes.NewBuffer(bookingBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	var bookingResp map[string]map[string]string
	json.NewDecoder(w.Body).Decode(&bookingResp)
	bookingID := bookingResp["booking"]["id"]

	req = httptest.NewRequest("POST", "/bookings/"+bookingID+"/cancel", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.SetPathValue("bookingId", bookingID)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for cancel, got %d: %s", w.Code, w.Body.String())
	}

	var cancelResp map[string]map[string]string
	json.NewDecoder(w.Body).Decode(&cancelResp)

	if cancelResp["booking"]["status"] != "cancelled" {
		t.Fatalf("expected cancelled, got %s", cancelResp["booking"]["status"])
	}
}
