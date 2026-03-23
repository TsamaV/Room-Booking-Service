package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string
	Role         string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
