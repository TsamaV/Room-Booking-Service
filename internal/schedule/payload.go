package schedule

type CreateScheduleRequest struct {
	DayOfWeek []int64 `json:"day_of_week"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
}

type ScheduleResponse struct {
	ID        string  `json:"id"`
	RoomID    string  `json:"room_id"`
	DayOfWeek []int64 `json:"day_of_week"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
}