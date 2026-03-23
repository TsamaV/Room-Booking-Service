package booking

import (
	"Avito/internal/slot"
	"Avito/internal/user"
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SlotID         uuid.UUID `gorm:"type:uuid;not null"`
	Slot           slot.Slot      `gorm:"foreignKey:SlotID"`
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	User           user.User      `gorm:"foreignKey:UserID"`
	Status         string    `gorm:"not null;default:active"`
	ConferenceLink string
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}
