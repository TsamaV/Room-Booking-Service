package schedule_test

import (
	"Avito/internal/schedule"
	"Avito/internal/slot"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type MockScheduleRepo struct {
	schedule *schedule.Schedule
}

func (m *MockScheduleRepo) Create(s *schedule.Schedule) (*schedule.Schedule, error) {
	s.ID = uuid.New()
	return s, nil
}

func (m *MockScheduleRepo) GetByRoomID(roomID uuid.UUID) (*schedule.Schedule, error) {
	if m.schedule != nil {
		return m.schedule, nil
	}
	return nil, nil
}

type MockSlotRepo struct {
	slots []slot.Slot
}

func (m *MockSlotRepo) BulkCreate(slots []slot.Slot) error {
	m.slots = slots
	return nil
}
func (m *MockSlotRepo) GetByID(id uuid.UUID) (*slot.Slot, error) { return nil, nil }
func (m *MockSlotRepo) Create(s *slot.Slot) (*slot.Slot, error)  { return s, nil }
func (m *MockSlotRepo) GetByRoomAndDate(roomID uuid.UUID, date time.Time) ([]slot.Slot, error) {
	return nil, nil
}

func TestCreateScheduleSuccess(t *testing.T) {
	scheduleRepo := &MockScheduleRepo{}
	slotRepo := &MockSlotRepo{}

	service := schedule.NewScheduleService(scheduleRepo, slotRepo, &MockBookingRepo{})

	startTime, _ := time.Parse("15:04", "09:00")
	endTime, _ := time.Parse("15:04", "18:00")

	s, err := service.Create(uuid.New(), []int64{1, 2, 3, 4, 5}, startTime, endTime)
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("expected schedule, got nil")
	}
}

func TestCreateScheduleAlreadyExists(t *testing.T) {
	scheduleRepo := &MockScheduleRepo{
		schedule: &schedule.Schedule{
			ID:        uuid.New(),
			DayOfWeek: pq.Int64Array{1, 2, 3},
		},
	}
	slotRepo := &MockSlotRepo{}

	service := schedule.NewScheduleService(scheduleRepo, slotRepo, &MockBookingRepo{})

	startTime, _ := time.Parse("15:04", "09:00")
	endTime, _ := time.Parse("15:04", "18:00")

	_, err := service.Create(uuid.New(), []int64{1, 2, 3}, startTime, endTime)
	if err == nil {
		t.Fatal("expected error for existing schedule, got nil")
	}
}

func TestGenerateSlotsCount(t *testing.T) {
	scheduleRepo := &MockScheduleRepo{}
	slotRepo := &MockSlotRepo{}

	service := schedule.NewScheduleService(scheduleRepo, slotRepo, &MockBookingRepo{})

	startTime, _ := time.Parse("15:04", "09:00")
	endTime, _ := time.Parse("15:04", "10:00")

	_, err := service.Create(uuid.New(), []int64{0, 1, 2, 3, 4, 5, 6}, startTime, endTime)
	if err != nil {
		t.Fatal(err)
	}

	if len(slotRepo.slots) == 0 {
		t.Fatal("expected slots to be created")
	}
}
