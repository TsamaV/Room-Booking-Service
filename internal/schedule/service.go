package schedule

import (
	"Avito/internal/slot"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ScheduleRepo interface {
	Create(s *Schedule) (*Schedule, error)
	GetByRoomID(roomID uuid.UUID) (*Schedule, error)
}

type SlotRepo interface {
	BulkCreate(slots []slot.Slot) error
	GetByID(id uuid.UUID) (*slot.Slot, error)
	GetByRoomAndDate(roomID uuid.UUID, date time.Time) ([]slot.Slot, error)
}

type BookingRepo interface {
	IsSlotBooked(slotID uuid.UUID) (bool, error)
}

type ScheduleService struct {
	scheduleRepo ScheduleRepo
	slotRepo     SlotRepo
	bookingRepo  BookingRepo
}

func NewScheduleService(scheduleRepo ScheduleRepo, slotRepo SlotRepo, bookingRepo BookingRepo) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		slotRepo:     slotRepo,
		bookingRepo:  bookingRepo,
	}
}

func (service *ScheduleService) Create(roomID uuid.UUID, dayOfWeek []int64, startTime, endTime time.Time) (*Schedule, error) {
	existing, _ := service.scheduleRepo.GetByRoomID(roomID)
	if existing != nil {
		return nil, errors.New("schedule already exists for this room")
	}

	newSchedule := Schedule{
		RoomID:    roomID,
		DayOfWeek: pq.Int64Array(dayOfWeek),
		StartTime: startTime,
		EndTime:   endTime,
	}
	savedSchedule, err := service.scheduleRepo.Create(&newSchedule)
	if err != nil {
		return nil, err
	}

	slots := generateSlots(savedSchedule)
	if len(slots) > 0 {
		if err := service.slotRepo.BulkCreate(slots); err != nil {
			return nil, err
		}
	}
	return savedSchedule, nil
}

func (service *ScheduleService) GetSlots(roomID uuid.UUID, date time.Time) ([]slot.Slot, error) {
	slots, err := service.slotRepo.GetByRoomAndDate(roomID, date)
	if err != nil {
		return nil, err
	}

	var freeSlots []slot.Slot
	for _, s := range slots {
		booked, _ := service.bookingRepo.IsSlotBooked(s.ID)
		if !booked {
			freeSlots = append(freeSlots, s)
		}
	}

	return freeSlots, nil
}

func generateSlots(schedule *Schedule) []slot.Slot {
	var slots []slot.Slot
	now := time.Now().UTC()

	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, i)

		if !containsDay(schedule.DayOfWeek, int64(date.Weekday())) {
			continue
		}

		slotStart := time.Date(date.Year(), date.Month(), date.Day(), schedule.StartTime.Hour(), schedule.StartTime.Minute(), 0, 0, time.UTC)
		endOfDay := time.Date(date.Year(), date.Month(), date.Day(), schedule.EndTime.Hour(), schedule.EndTime.Minute(), 0, 0, time.UTC)

		for slotStart.Before(endOfDay) {
			slotEnd := slotStart.Add(30 * time.Minute)
			if slotEnd.After(endOfDay) {
				break
			}

			slots = append(slots, slot.Slot{
				RoomID:     schedule.RoomID,
				ScheduleID: schedule.ID,
				StartTime:  slotStart,
				EndTime:    slotEnd,
			})

			slotStart = slotEnd
		}
	}
	return slots
}

func toGoWeekday(d int64) int64 {
	if d == 7 {
		return 0
	}
	return d
}

func containsDay(days pq.Int64Array, day int64) bool {
	for _, d := range days {
		if toGoWeekday(d) == day {
			return true
		}
	}
	return false
}