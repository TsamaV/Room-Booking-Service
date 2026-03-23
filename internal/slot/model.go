package slot

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoomID     uuid.UUID `gorm:"type:uuid;not null"`
	ScheduleID uuid.UUID `gorm:"type:uuid;not null"`
	StartTime  time.Time `gorm:"not null"`
	EndTime    time.Time `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
