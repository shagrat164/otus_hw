package sqlstorage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/models"
)

// Storage - SQL реализация хранилища событий с использованием pgx.
type Storage struct {
	db  *pgxpool.Pool
	dsn string // "postgres://user:password@host:port/dbname?sslmode=disable"
}

// New создаёт новый экземпляр SQL хранилища.
func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

// Connect подключается к базе данных.
func (s *Storage) Connect(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, s.dsn)
	if err != nil {
		return err
	}
	s.db = pool
	return nil
}

// Close закрывает подключение к базе данных.
func (s *Storage) Close() {
	s.db.Close()
}

// CreateEvent создаёт новое событие.
func (s *Storage) CreateEvent(event models.Event) error {
	query := `
	INSERT INTO events (id, title, start_time, end_time, description, user_id, reminder) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.Exec(context.Background(), query,
		event.ID,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.UserID,
		event.Reminder,
	)
	if err != nil {
		return err
	}

	return nil
}

// UpdateEvent обновляет существующее событие по id.
func (s *Storage) UpdateEvent(id uuid.UUID, event models.Event) error {
	query := `
	UPDATE events
	SET title = $1, start_time = $2, end_time = $3, description = $4, user_id = $5, reminder = $6
	WHERE id = $7
	`

	_, err := s.db.Exec(context.Background(), query,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.UserID,
		event.Reminder,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteEvent удаляет событие по id.
func (s *Storage) DeleteEvent(id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`

	_, err := s.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}

// GetEventsForDay возвращает список событий на указанный день.
func (s *Storage) GetEventsForDay(date time.Time) ([]models.Event, error) {
	query := `
	SELECT id, title, start_time, end_time, description, user_id, reminder
	FROM events
	WHERE start_time::date = $1::date
	`

	rows, err := s.db.Query(context.Background(), query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.StartTime,
			&event.EndTime,
			&event.Description,
			&event.UserID,
			&event.Reminder,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, event)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

// GetEventsForWeek возвращает список событий на неделю, начиная с указанной даты.
func (s *Storage) GetEventsForWeek(startDate time.Time) ([]models.Event, error) { //nolint:dupl
	endDate := startDate.AddDate(0, 0, 7)
	query := `
	SELECT id, title, start_time, end_time, description, user_id, reminder
	FROM events
	WHERE start_time >= $1 AND start_time < $2
	`

	rows, err := s.db.Query(context.Background(), query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.StartTime,
			&event.EndTime,
			&event.Description,
			&event.UserID,
			&event.Reminder,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, event)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

// GetEventsForMonth возвращает список событий на месяц, начиная с указанной даты.
func (s *Storage) GetEventsForMonth(startDate time.Time) ([]models.Event, error) { //nolint:dupl
	endDate := startDate.AddDate(0, 1, 0)
	query := `
	SELECT id, title, start_time, end_time, description, user_id, reminder
	FROM events
	WHERE start_time >= $1 AND start_time < $2
	`

	rows, err := s.db.Query(context.Background(), query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.StartTime,
			&event.EndTime,
			&event.Description,
			&event.UserID,
			&event.Reminder,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, event)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}
