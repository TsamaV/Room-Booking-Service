// @title Room Booking Service
// @version 1.0
// @description Сервис бронирования переговорок
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"Avito/configs"
	_ "Avito/docs"
	"Avito/internal/auth"
	"Avito/internal/booking"
	"Avito/internal/db"
	"Avito/internal/room"
	"Avito/internal/schedule"
	"Avito/internal/slot"
	"Avito/internal/user"
	"fmt"
	"github.com/swaggo/http-swagger"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	database.AutoMigrate(&user.User{}, &room.Room{}, &schedule.Schedule{}, &slot.Slot{}, &booking.Booking{})
	database.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_active_booking_slot ON bookings (slot_id) WHERE status = 'active'`)
	router := http.NewServeMux()

	// Repositories
	bookingRepository := booking.NewBookingRepository(database)
	roomRepository := room.NewRoomRepository(database)
	scheduleRepository := schedule.NewScheduleRepository(database)
	slotRepository := slot.NewSlotRepository(database)
	userRepository := user.NewUserRepository(database)

	// Services
	scheduleService := schedule.NewScheduleService(scheduleRepository, slotRepository, bookingRepository)
	bookingService := booking.NewBookingService(bookingRepository, slotRepository)
	authService := auth.NewAuthService(userRepository, conf.Auth.Secret)
	roomService := room.NewRoomService(roomRepository)

	// Handler
	auth.NewAuthHandler(router, authService)
	room.NewRoomHandler(router, roomService, conf.Auth.Secret)
	schedule.NewScheduleHandler(router, scheduleService, conf.Auth.Secret)
	booking.NewBookingHandler(router, bookingService, conf.Auth.Secret)

	router.HandleFunc("GET /_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})
	router.Handle("GET /swagger/", httpSwagger.WrapHandler)

	return router
}

func main() {
	app := App()

	server := http.Server{
		Addr:    ":8080",
		Handler: app,
	}

	fmt.Println("Server is listening on 8080")
	server.ListenAndServe()
}
