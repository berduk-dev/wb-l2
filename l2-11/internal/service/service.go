package service

import (
	"calendar-server/internal/model"
	"calendar-server/internal/repo"
	"context"
)

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
