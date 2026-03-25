package schedule

import (
	"Avito/internal/db"

	"github.com/google/uuid"
)

type ScheduleRepository struct {
	DataBase *db.Db
}

func NewScheduleRepository(database *db.Db) *ScheduleRepository {
	return &ScheduleRepository{
		DataBase:database,
	}
}

func (repo *ScheduleRepository) Create(schedule *Schedule) (*Schedule, error) {
	result := repo.DataBase.Create(schedule) 
	if result.Error != nil {
		return nil, result.Error
	}
	return schedule, nil
} 

func (repo *ScheduleRepository) GetByRoomID(roomID uuid.UUID) (*Schedule, error) {
	var schedule Schedule
	result := repo.DataBase.Where("room_id = ?", roomID).First(&schedule)
	if result.Error != nil {
		return nil, result.Error
	}
	return &schedule, nil
}