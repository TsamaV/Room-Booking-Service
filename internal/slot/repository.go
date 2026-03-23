package slot

import (
	"Avito/pkg/db"
	"time"

	"github.com/google/uuid"
)

type SlotRepository struct {
	DataBase *db.Db
}

func NewSlotRepository(database *db.Db) *SlotRepository {
	return &SlotRepository{
		DataBase: database,
	}
}

func (repo *SlotRepository) Create(slot *Slot) (*Slot, error) {
	result := repo.DataBase.Create(slot)
	if result.Error != nil {
		return nil, result.Error
	}
	return slot, nil
}

func (repo *SlotRepository) GetByID(id uuid.UUID) (*Slot, error) {
	var slot Slot
	result := repo.DataBase.Where("id = ?", id).First(&slot)
	if result.Error != nil {
		return nil, result.Error
	}
	return &slot, nil
}

func (repo *SlotRepository) BulkCreate(slots []Slot) error {
	result := repo.DataBase.Create(slots)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *SlotRepository) GetByRoomAndDate(roomID uuid.UUID, date time.Time) ([]Slot, error) {
	var slots []Slot
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	result := repo.DataBase.Where(
		"room_id = ? AND start_time >= ? AND start_time < ?",
		roomID, startOfDay, endOfDay,
	).Order("start_time").Find(&slots)

	if result.Error != nil {
		return nil, result.Error
	}
	return slots, nil
}