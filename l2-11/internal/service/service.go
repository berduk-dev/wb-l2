package service

import (
	"calendar-server/internal/model"
	"calendar-server/internal/repo"
	"context"
	"log"
	"time"
)

const layout = "2006-01-02"

type Service struct {
	repo repo.Repository
}

func New(repo repo.Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s *Service) CreateEvent(ctx context.Context, req model.CreateEventRequest) (int, error) {
	id, err := s.repo.CreateEvent(ctx, req)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) UpdateEvent(ctx context.Context, event model.UpdateEventRequest) error {
	err := s.repo.UpdateEvent(ctx, event)
	return err
}

func (s *Service) DeleteEvent(ctx context.Context, eventID int) error {
	err := s.repo.DeleteEvent(ctx, eventID)
	return err
}

func (s *Service) EventsForToday(ctx context.Context, userID int) ([]model.Event, error) {
	events, err := s.repo.EventsForDay(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var dailyEvents []model.Event
	for _, e := range events {
		eventDate, err := time.Parse(layout, e.Date)
		if err != nil {
			log.Println("parse error:", e.Date, err)
			continue
		}

		eventDate = eventDate.Truncate(24 * time.Hour)

		if eventDate.Year() == today.Year() &&
			eventDate.Month() == today.Month() &&
			eventDate.Day() == today.Day() {
			dailyEvents = append(dailyEvents, e)
		}
	}

	return dailyEvents, nil
}
