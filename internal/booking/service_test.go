package booking_test

import (
	"Avito/internal/booking"
	"Avito/internal/slot"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type MockBookingRepo struct{}

func (m *MockBookingRepo) Create(b *booking.Booking) (*booking.Booking, error) {
	return b, nil
}
func (m *MockBookingRepo) GetByID(id uuid.UUID) (*booking.Booking, error) {
	return nil, errors.New("not found")
}
func (m *MockBookingRepo) GetActiveBySlotID(slotID uuid.UUID) (*booking.Booking, error) {
	return nil, errors.New("not found")
}
func (m *MockBookingRepo) Cancel(id uuid.UUID) (*booking.Booking, error) {
	return &booking.Booking{Status: "cancelled"}, nil
}
func (m *MockBookingRepo) GetMyBookings(userID uuid.UUID) ([]booking.Booking, error) {
	return nil, nil
}
func (m *MockBookingRepo) GetAllPaginated(page, limit int) ([]booking.Booking, int64, error) {
	return nil, 0, nil
}

type MockBookingRepoOccupied struct{}

func (m *MockBookingRepoOccupied) Create(b *booking.Booking) (*booking.Booking, error) {
	return b, nil
}
func (m *MockBookingRepoOccupied) GetByID(id uuid.UUID) (*booking.Booking, error) {
	return nil, errors.New("not found")
}
func (m *MockBookingRepoOccupied) GetActiveBySlotID(slotID uuid.UUID) (*booking.Booking, error) {
	return &booking.Booking{Status: "active"}, nil
}
func (m *MockBookingRepoOccupied) Cancel(id uuid.UUID) (*booking.Booking, error) {
	return &booking.Booking{Status: "cancelled"}, nil
}
func (m *MockBookingRepoOccupied) GetMyBookings(userID uuid.UUID) ([]booking.Booking, error) {
	return nil, nil
}
func (m *MockBookingRepoOccupied) GetAllPaginated(page, limit int) ([]booking.Booking, int64, error) {
	return nil, 0, nil
}

type MockBookingRepoWithBooking struct {
	booking *booking.Booking
}

func (m *MockBookingRepoWithBooking) Create(b *booking.Booking) (*booking.Booking, error) {
	return b, nil
}
func (m *MockBookingRepoWithBooking) GetByID(id uuid.UUID) (*booking.Booking, error) {
	return m.booking, nil
}
func (m *MockBookingRepoWithBooking) GetActiveBySlotID(slotID uuid.UUID) (*booking.Booking, error) {
	return nil, errors.New("not found")
}
func (m *MockBookingRepoWithBooking) Cancel(id uuid.UUID) (*booking.Booking, error) {
	m.booking.Status = "cancelled"
	return m.booking, nil
}
func (m *MockBookingRepoWithBooking) GetMyBookings(userID uuid.UUID) ([]booking.Booking, error) {
	return []booking.Booking{*m.booking}, nil
}
func (m *MockBookingRepoWithBooking) GetAllPaginated(page, limit int) ([]booking.Booking, int64, error) {
	return []booking.Booking{*m.booking}, 1, nil
}

type MockSlotRepo struct{}

func (m *MockSlotRepo) GetByID(id uuid.UUID) (*slot.Slot, error) {
	return &slot.Slot{
		ID:        id,
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}, nil
}

type MockSlotRepoPast struct{}

func (m *MockSlotRepoPast) GetByID(id uuid.UUID) (*slot.Slot, error) {
	return &slot.Slot{
		ID:        id,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now().Add(-30 * time.Minute),
	}, nil
}

func TestCreateBookingSuccess(t *testing.T) {
	service := booking.NewBookingService(&MockBookingRepo{}, &MockSlotRepo{})

	b, err := service.Create(uuid.New(), uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if b.Status != "active" {
		t.Fatalf("expected status active, got %s", b.Status)
	}
}

func TestCreateBookingPastSlot(t *testing.T) {
	service := booking.NewBookingService(&MockBookingRepo{}, &MockSlotRepoPast{})

	_, err := service.Create(uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for past slot, got nil")
	}
}

func TestCreateBookingAlreadyBooked(t *testing.T) {
	service := booking.NewBookingService(&MockBookingRepoOccupied{}, &MockSlotRepo{})

	_, err := service.Create(uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for occupied slot, got nil")
	}
}

func TestCancelBookingSuccess(t *testing.T) {
	userID := uuid.New()

	repo := &MockBookingRepoWithBooking{
		booking: &booking.Booking{
			ID:     uuid.New(),
			UserID: userID,
			Status: "active",
		},
	}

	service := booking.NewBookingService(repo, &MockSlotRepo{})

	b, err := service.Cancel(uuid.New(), userID)
	if err != nil {
		t.Fatal(err)
	}
	if b.Status != "cancelled" {
		t.Fatalf("expected cancelled, got %s", b.Status)
	}
}

func TestCancelBookingWrongUser(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	repo := &MockBookingRepoWithBooking{
		booking: &booking.Booking{
			ID:     uuid.New(),
			UserID: otherUserID,
			Status: "active",
		},
	}

	service := booking.NewBookingService(repo, &MockSlotRepo{})

	_, err := service.Cancel(uuid.New(), userID)
	if err == nil {
		t.Fatal("expected error for wrong user, got nil")
	}
}

func TestCancelBookingIdempotent(t *testing.T) {
	userID := uuid.New()

	repo := &MockBookingRepoWithBooking{
		booking: &booking.Booking{
			ID:     uuid.New(),
			UserID: userID,
			Status: "cancelled",
		},
	}

	service := booking.NewBookingService(repo, &MockSlotRepo{})

	b, err := service.Cancel(uuid.New(), userID)
	if err != nil {
		t.Fatal(err)
	}
	if b.Status != "cancelled" {
		t.Fatalf("expected cancelled, got %s", b.Status)
	}
}

type MockSlotRepoNotFound struct{}

func (m *MockSlotRepoNotFound) GetByID(id uuid.UUID) (*slot.Slot, error) {
	return nil, errors.New("not found")
}

func TestCreateBookingSlotNotFound(t *testing.T) {
	service := booking.NewBookingService(&MockBookingRepo{}, &MockSlotRepoNotFound{})

	_, err := service.Create(uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected error for not found slot, got nil")
	}
}

func TestGetMyBookings(t *testing.T) {
	userID := uuid.New()

	repo := &MockBookingRepoWithBooking{
		booking: &booking.Booking{
			ID:     uuid.New(),
			UserID: userID,
			Status: "active",
		},
	}

	service := booking.NewBookingService(repo, &MockSlotRepo{})

	bookings, err := service.GetMyBookings(userID)
	if err != nil {
		t.Fatal(err)
	}
	if len(bookings) == 0 {
		t.Fatal("expected bookings, got empty")
	}
}

func TestGetAllPaginated(t *testing.T) {
	userID := uuid.New()

	repo := &MockBookingRepoWithBooking{
		booking: &booking.Booking{
			ID:     uuid.New(),
			UserID: userID,
			Status: "active",
		},
	}

	service := booking.NewBookingService(repo, &MockSlotRepo{})

	bookings, total, err := service.GetAllPaginated(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total == 0 {
		t.Fatal("expected total > 0")
	}
	if len(bookings) == 0 {
		t.Fatal("expected bookings, got empty")
	}
}
