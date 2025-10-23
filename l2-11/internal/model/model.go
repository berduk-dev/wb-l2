package model

import "time"

type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
}

type CreateEventRequest struct {
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
}
