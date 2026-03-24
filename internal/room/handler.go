package room

import (
	"Avito/pkg/middleware"
	"Avito/pkg/req"
	"Avito/pkg/res"
	"net/http"

	"github.com/google/uuid"
)

type RoomHandler struct {
	RoomService *RoomService
}

func NewRoomHandler(router *http.ServeMux, service *RoomService, jwtSecret string) {
	handler := &RoomHandler{
		RoomService: service,
	}

	router.Handle("POST /rooms/create", middleware.IsAuthed(middleware.AdminOnly(http.HandlerFunc(handler.Create)), jwtSecret))
	router.Handle("GET /rooms/list", middleware.IsAuthed(http.HandlerFunc(handler.GetAll), jwtSecret))
}

func NewRoomHandler2(service *RoomService) *RoomHandler {
	return &RoomHandler{RoomService: service}
}

// @Summary Создать комнату
// @Tags Rooms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateRoomRequest true "Данные комнаты"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /rooms/create [post]
func (handler *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.Decode[CreateRoomRequest](r)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	room, err := handler.RoomService.Create(body.Name, body.Description, body.Capacity)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INTERNAL_ERROR", err.Error())
		return
	}

	res.JSON(w, http.StatusCreated, map[string]any{
		"room": RoomResponse{
			ID:          room.ID.String(),
			Name:        room.Name,
			Description: room.Description,
			Capacity:    room.Capacity,
			CreatedAt:   room.CreatedAt.String(),
		},
	})
}

// @Summary Список всех комнат
// @Tags Rooms
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /rooms/list [get]
func (handler *RoomHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rooms, err := handler.RoomService.GetAll()
	if err != nil {
		res.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	var response []RoomResponse
	for _, room := range rooms {
		response = append(response, RoomResponse{
			ID:          room.ID.String(),
			Name:        room.Name,
			Description: room.Description,
			Capacity:    room.Capacity,
			CreatedAt:   room.CreatedAt.String(),
		})
	}

	res.JSON(w, http.StatusOK, map[string]any{
		"rooms": response,
	})
}

func (handler *RoomHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		res.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid room id")
		return
	}

	room, err := handler.RoomService.GetByID(id)
	if err != nil {
		res.Error(w, http.StatusNotFound, "ROOM_NOT_FOUND", "room not found")
		return
	}

	res.JSON(w, http.StatusOK, map[string]any{
		"room": RoomResponse{
			ID:          room.ID.String(),
			Name:        room.Name,
			Description: room.Description,
			Capacity:    room.Capacity,
			CreatedAt:   room.CreatedAt.String(),
		},
	})
}