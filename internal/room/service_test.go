package room_test

import (
	"Avito/internal/room"
	"testing"

	"github.com/google/uuid"
)

type MockRoomRepo struct{}

func (m *MockRoomRepo) Create(r *room.Room) (*room.Room, error) {
	return r, nil
}
func (m *MockRoomRepo) GetAll() ([]room.Room, error) {
	return []room.Room{{ID: uuid.New(), Name: "Test"}}, nil
}
func (m *MockRoomRepo) GetByID(id uuid.UUID) (*room.Room, error) {
	return &room.Room{ID: id, Name: "Test"}, nil
}

func TestRoomCreate(t *testing.T) {
	service := room.NewRoomService(&MockRoomRepo{})
	r, err := service.Create("Test", "Desc", 10)
	if err != nil {
		t.Fatal(err)
	}
	if r.Name != "Test" {
		t.Fatalf("expected Test, got %s", r.Name)
	}
}

func TestRoomGetAll(t *testing.T) {
	service := room.NewRoomService(&MockRoomRepo{})
	rooms, err := service.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rooms) == 0 {
		t.Fatal("expected rooms, got empty")
	}
}

func TestRoomGetByID(t *testing.T) {
	service := room.NewRoomService(&MockRoomRepo{})
	id := uuid.New()
	r, err := service.GetByID(id)
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != id {
		t.Fatalf("expected %s, got %s", id, r.ID)
	}
}
