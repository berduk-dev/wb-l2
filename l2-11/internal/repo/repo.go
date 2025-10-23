package repo

import (
	"calendar-server/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) Repository {
	return Repository{
		db: db,
	}
}

func (r *Repository) CreateEvent(ctx context.Context, req model.CreateEventRequest) (int, error) {
	var id int

	err := r.db.QueryRow(ctx, `
		INSERT INTO events
					(user_id, title, date)
					VALUES ($1, $2, $3) RETURNING id`,
		req.UserID, req.Title, req.Date,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error CreateEvent INSERT")
	}

	return id, nil
}
