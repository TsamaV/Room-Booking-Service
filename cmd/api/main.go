package main

import (
	"Avito/configs"
	"Avito/internal/auth"
	"Avito/internal/booking"
	"Avito/internal/room"
	"Avito/internal/schedule"
	"Avito/internal/slot"
	"Avito/internal/user"
	"Avito/pkg/db"
	"fmt"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	bookingRepository := booking.NewBookingRepository(db)
	roomRepository := room.NewRoomRepository(db)
	scheduleRepository := schedule.NewScheduleRepository(db)
	slotRepository := slot.NewSlotRepository(db)
	userRepository := user.NewUserRepository(db)

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
