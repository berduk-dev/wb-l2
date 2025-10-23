package model

type Event struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
}

type CreateEventRequest struct {
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
}

type DeleteEventRequest struct {
	ID int `json:"id"`
}

type UpdateEventRequest struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}
