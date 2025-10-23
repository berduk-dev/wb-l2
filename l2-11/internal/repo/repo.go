package repo

import (
	"calendar-server/internal/errs"
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

func (r *Repository) UpdateEvent(ctx context.Context, event model.UpdateEventRequest) error {
	cmdTag, err := r.db.Exec(ctx, `
		UPDATE events
		SET title = $1, date = $2
		WHERE id = $3 `,
		event.Title, event.Date, event.ID)
	if err != nil {
		return fmt.Errorf("error UpdateEvent Exec: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errs.ErrEventNotFound
	}

	return nil
}

func (r *Repository) DeleteEvent(ctx context.Context, eventID int) error {
	cmdTag, err := r.db.Exec(ctx, `
		DELETE FROM events
		WHERE id = $1`,
		eventID)
	if err != nil {
		return fmt.Errorf("error DeleteEvent Exec: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errs.ErrEventNotFound
	}
	return nil
}

func (r *Repository) EventsForDay(ctx context.Context, userID int) ([]model.Event, error) {
	rows, err := r.db.Query(ctx, `
			SELECT * FROM events
			WHERE user_id = $1`,
		userID,
	)

	var events []model.Event
	for rows.Next() {
		var event model.Event
		err = rows.Scan(&event.ID, &event.UserID, &event.Title, &event.Date)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
