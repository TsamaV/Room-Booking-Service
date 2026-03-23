package schedule

import (
	"Avito/internal/room"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Schedule struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoomID    uuid.UUID     `gorm:"type:uuid;not null"`
	Room      room.Room          `gorm:"foreignKey:RoomID"`
	DayOfWeek pq.Int64Array `gorm:"type:integer[]"`
	StartTime time.Time     `gorm:"not null"`
	EndTime   time.Time     `gorm:"not null"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
}
