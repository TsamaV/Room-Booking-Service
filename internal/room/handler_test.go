package room_test

import (
	"Avito/internal/room"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestRoomHandlerCreate(t *testing.T) {
	service := room.NewRoomService(&MockRoomRepo{})
	handler := room.NewRoomHandler2(service)

	body, _ := json.Marshal(map[string]any{
		"name":        "Test Room",
		"description": "Test",
		"capacity":    10,
	})

	req := httptest.NewRequest("POST", "/rooms/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestRoomHandlerGetAll(t *testing.T) {
	service := room.NewRoomService(&MockRoomRepo{})
	handler := room.NewRoomHandler2(service)

	req := httptest.NewRequest("GET", "/rooms/list", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRoomHandlerGetByID(t *testing.T) {
	service := room.NewRoomService(&MockRoomRepo{})
	handler := room.NewRoomHandler2(service)

	id := uuid.New()
	req := httptest.NewRequest("GET", "/rooms/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}