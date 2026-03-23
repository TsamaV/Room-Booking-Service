package booking

type CreateBookingRequest struct {
	SlotID string `json:"slot_id"`
}

type BookingResponse struct {
	ID             string `json:"id"`
	SlotID         string `json:"slot_id"`
	UserID         string `json:"user_id"`
	Status         string `json:"status"`
	ConferenceLink string `json:"conference_link,omitempty"`
	CreatedAt      string `json:"created_at"`
}

type PaginatedBookingsResponse struct {
	Data  []BookingResponse `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}