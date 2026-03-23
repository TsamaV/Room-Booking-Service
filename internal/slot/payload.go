package slot

type SlotResponse struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	IsBooked  bool   `json:"is_booked"`
}