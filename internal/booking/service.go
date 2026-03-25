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

	booking := &Booking{
		SlotID: slotID,
		UserID: userID,
		Status: "active",
	}

	created, err := service.bookingRepo.Create(booking)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, errors.New("slot is already booked")
		}
		return nil, err
	}
	return created, nil
}

func isDuplicateKeyError(err error) bool {
	return err != nil && (err.Error() == "ERROR: duplicate key value violates unique constraint \"idx_active_booking_slot\" (SQLSTATE 23505)" ||
		err.Error() == "duplicate key value violates unique constraint \"idx_active_booking_slot\"")
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
