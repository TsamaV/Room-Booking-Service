package main

import (
	"Avito/internal/booking"
	"Avito/internal/room"
	"Avito/internal/schedule"
	"Avito/internal/slot"
	"Avito/internal/user"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&user.User{}, &room.Room{}, &schedule.Schedule{}, &slot.Slot{}, &booking.Booking{})
}