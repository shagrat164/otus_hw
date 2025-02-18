package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/models"
)

type App struct { // TODO
	config  config.Config
	storage Storage
	logger  Logger
}

// Logger - интерфейс для работы с логами.
type Logger interface {
	// Debug выводит сообщение уровня DEBUG.
	Debug(msg string)

	// Info выводит сообщение уровня INFO.
	Info(msg string)

	// Warn выводит сообщение уровня WARN.
	Warn(msg string)

	// Error выводит сообщение уровня ERROR.
	Error(msg string)
}

// Storage - интерфейс для работы с хранилищем событий.
type Storage interface {
	// CreateEvent создаёт новое событие.
	CreateEvent(event models.Event) error

	// UpdateEvent обновляет существующее событие по id.
	UpdateEvent(id uuid.UUID, event models.Event) error

	// DeleteEvent удаляет событие по id.
	DeleteEvent(id uuid.UUID) error

	// GetEventsForDay возвращает список событий на указанный день.
	GetEventsForDay(date time.Time) ([]models.Event, error)

	// GetEventsForWeek возвращает список событий на неделю, начиная с указанной даты.
	GetEventsForWeek(startDate time.Time) ([]models.Event, error)

	// GetEventsForMonth возвращает список событий на месяц, начиная с указанной даты.
	GetEventsForMonth(startDate time.Time) ([]models.Event, error)
}

func New(logger Logger, storage Storage, config config.Config) *App {
	return &App{
		config:  config,
		storage: storage,
		logger:  logger,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error { //nolint:revive
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
