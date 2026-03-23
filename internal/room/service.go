package room

import "github.com/google/uuid"

type RoomRepo interface {
	Create(r *Room) (*Room, error)
	GetAll() ([]Room, error)
	GetByID(id uuid.UUID) (*Room, error)
}

type RoomService struct {
	repo RoomRepo
}

func NewRoomService(repo RoomRepo) *RoomService {
	return &RoomService{
		repo: repo,
	}
}

func (service *RoomService) Create(name, description string, capacity int) (*Room, error) {
	room := &Room{
		Name:        name,
		Description: description,
		Capacity:    capacity,
	}
	return service.repo.Create(room)
}

func (service *RoomService) GetAll() ([]Room, error) {
	return service.repo.GetAll()
}

func (s *RoomService) GetByID(id uuid.UUID) (*Room, error) {
	return s.repo.GetByID(id)
}
