package booking

import (
	"Avito/pkg/db"
	"time"

	"github.com/google/uuid"
)

type BookingRepository struct {
	DataBase *db.Db
}

func NewBookingRepository(database *db.Db) *BookingRepository {
	return &BookingRepository{
		DataBase: database,
	}
}

func (repo *BookingRepository) Create(booking *Booking) (*Booking, error) {
	result := repo.DataBase.Create(booking)
	if result.Error != nil {
		return nil, result.Error
	}
	return booking, nil
}

func (repo *BookingRepository) GetByID(id uuid.UUID) (*Booking, error) {
	var booking Booking
	result := repo.DataBase.Where("id = ?", id).First(&booking)
	if result.Error != nil {
		return nil, result.Error
	}
	return &booking, nil
}

func (repo *BookingRepository) Cancel(id uuid.UUID) (*Booking, error) {
	result := repo.DataBase.Model(&Booking{}).Where("id = ? ", id).Update("status", "cancelled")
	if result.Error != nil {
		return nil, result.Error
	}
	return repo.GetByID(id)
}

func (repo *BookingRepository) GetAllPaginated(page, limit int) ([]Booking, int64, error) {
	var bookings []Booking
	var total int64

	repo.DataBase.Model(Booking{}).Count(&total)
	result := repo.DataBase.Offset((page - 1) * limit).Limit(limit).Find(&bookings)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return bookings, total, nil
}

func (repo *BookingRepository) GetMyBookings(userID uuid.UUID) ([]Booking, error) {
	var bookings []Booking
	result := repo.DataBase.Joins("JOIN slots ON slots.id = bookings.slot_id").Where("bookings.user_id = ? AND bookings.status = ? AND slots.start_time > ?", userID, "active", time.Now().UTC()).Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}

func (repo *BookingRepository) GetActiveBySlotID(slotID uuid.UUID) (*Booking, error) {
	var booking Booking
	result := repo.DataBase.Where("slot_id = ? AND status = ?", slotID, "active").First(&booking)
	if result.Error != nil {
		return nil, result.Error
	}
	return &booking, nil
}

func (repo *BookingRepository) IsSlotBooked(slotID uuid.UUID) (bool, error) {
	var count int64
	result := repo.DataBase.Model(&Booking{}).
		Where("slot_id = ? AND status = ?", slotID, "active").
		Count(&count)
	return count > 0, result.Error
}