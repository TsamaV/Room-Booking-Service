package booking

import (
	"Avito/internal/slot"
	"errors"
	"time"

	"github.com/google/uuid"
)

type BookingRepo interface {
	Create(b *Booking) (*Booking, error)
	GetByID(id uuid.UUID) (*Booking, error)
	Cancel(id uuid.UUID) (*Booking, error)
	GetActiveBySlotID(slotID uuid.UUID) (*Booking, error)
	GetMyBookings(userID uuid.UUID) ([]Booking, error)
	GetAllPaginated(page, limit int) ([]Booking, int64, error)
}

type SlotRepo interface {
	GetByID(id uuid.UUID) (*slot.Slot, error)
}

type BookingService struct {
	bookingRepo BookingRepo
	slotRepo    SlotRepo
}

func NewBookingService(bookingRepo BookingRepo, slotRepo SlotRepo) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		slotRepo:    slotRepo,
	}
}

func (service *BookingService) Create(userID, slotID uuid.UUID) (*Booking, error) {
	sl, err := service.slotRepo.GetByID(slotID)
	if err != nil {
		return nil, err
	}
	if sl.StartTime.Before(time.Now().UTC()) {
		return nil, errors.New("cannot book a slot in the past")
	}

	existing, _ := service.bookingRepo.GetActiveBySlotID(slotID)
	if existing != nil {
		return nil, errors.New("slot is already booked")
	}

	booking := &Booking{
		SlotID: slotID,
		UserID: userID,
		Status: "active",
	}
	return service.bookingRepo.Create(booking)
}

func (service *BookingService) Cancel(bookingID, userID uuid.UUID) (*Booking, error) {
	booking, err := service.bookingRepo.GetByID(bookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	if booking.UserID != userID {
		return nil, errors.New("you can only cancel your own bookings")
	}

	if booking.Status == "cancelled" {
		return booking, nil
	}
	return service.bookingRepo.Cancel(bookingID)
}

func (service *BookingService) GetMyBookings(userID uuid.UUID) ([]Booking, error) {
	return service.bookingRepo.GetMyBookings(userID)
}

func (service *BookingService) GetAllPaginated(page, limit int) ([]Booking, int64, error) {
	return service.bookingRepo.GetAllPaginated(page, limit)
}