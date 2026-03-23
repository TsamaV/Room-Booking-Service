package room

import (
	"Avito/pkg/db"

	"github.com/google/uuid"
)

type RoomRepository struct {
	DataBase *db.Db
}

func NewRoomRepository(database *db.Db) *RoomRepository {
	return &RoomRepository{
		DataBase: database,
	}
}

func (repo *RoomRepository) Create(room *Room) (*Room, error) {
	result := repo.DataBase.Create(room)
	if result.Error != nil {
		return nil, result.Error
	}
	return room, nil
}

func (repo *RoomRepository) GetAll() ([]Room, error) {
	var rooms []Room
	result := repo.DataBase.Find(&rooms)
	if result.Error != nil {
		return nil, result.Error
	}
	return rooms, nil
}

func (repo *RoomRepository) GetByID(id uuid.UUID) (*Room, error) {
	var room Room
	result := repo.DataBase.Where("id = ?", id).First(&room)
	if result.Error != nil {
		return nil, result.Error
	}
	return &room, nil
}
