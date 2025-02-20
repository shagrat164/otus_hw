package memorystorage

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/models"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/storage"
)

// Storage - реализация in-memory хранилища событий.
type Storage struct {
	events map[uuid.UUID]models.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]models.Event),
	}
}

// CreateEvent создаёт новое событие.
func (s *Storage) CreateEvent(event models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; exists {
		return storage.ErrEventAlreadyExist
	}

	s.events[event.ID] = event
	return nil
}

// UpdateEvent обновляет существующее событие по id.
func (s *Storage) UpdateEvent(id uuid.UUID, event models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return storage.ErrEventNotFound
	}

	s.events[id] = event
	return nil
}

// DeleteEvent удаляет событие по id.
func (s *Storage) DeleteEvent(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return storage.ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

// GetEventsForDay возвращает список событий на указанный день.
func (s *Storage) GetEventsForDay(date time.Time) ([]models.Event, error) {
	// Блокировка на чтение
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Event
	for _, event := range s.events {
		// Пробежаться по всем событиям
		if event.StartTime.Year() == date.Year() &&
			event.StartTime.Month() == date.Month() &&
			event.StartTime.Day() == date.Day() {
			// Добавить только с одинаковой датой
			result = append(result, event)
		}
	}

	return result, nil
}

// GetEventsForWeek возвращает список событий на неделю, начиная с указанной даты.
func (s *Storage) GetEventsForWeek(startDate time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Event
	endDate := startDate.AddDate(0, 0, 7)
	for _, event := range s.events {
		// Пробежаться по всем событиям
		if (event.StartTime.Year() >= startDate.Year() &&
			event.StartTime.Month() >= startDate.Month() &&
			event.StartTime.Day() >= startDate.Day()) &&
			event.StartTime.Before(endDate) {
			// Добавить только с нужной датой в диапазоне [startDate..endDate)
			result = append(result, event)
		}
	}

	return result, nil
}

// GetEventsForMonth возвращает список событий на месяц, начиная с указанной даты.
func (s *Storage) GetEventsForMonth(startDate time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Event
	endDate := startDate.AddDate(0, 1, 0)
	for _, event := range s.events {
		// Пробежаться по всем событиям
		if (event.StartTime.Year() >= startDate.Year() &&
			event.StartTime.Month() >= startDate.Month() &&
			event.StartTime.Day() >= startDate.Day()) &&
			event.StartTime.Before(endDate) {
			// Добавить только с нужной датой в диапазоне [startDate..endDate)
			result = append(result, event)
		}
	}

	return result, nil
}
