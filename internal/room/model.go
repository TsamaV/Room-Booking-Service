package room

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	Description string
	Capacity    int
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
