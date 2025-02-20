package memorystorage

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/models"
	"github.com/stretchr/testify/require"
)

// TestCreateEvent проверяет создание события.
func TestCreateEvent(t *testing.T) {
	s := New()
	event := models.Event{
		ID:        uuid.New(),
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}

	// Создаем событие
	err := s.CreateEvent(event)
	require.NoError(t, err)

	// Проверяем, что событие добавлено
	events, err := s.GetEventsForDay(event.StartTime)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, event, events[0])
}

// TestUpdateEvent проверяет обновление события.
func TestUpdateEvent(t *testing.T) {
	s := New()
	event := models.Event{
		ID:        uuid.New(),
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}

	// Создаем событие
	err := s.CreateEvent(event)
	require.NoError(t, err)

	// Обновляем событие
	updatedEvent := event
	updatedEvent.Title = "Updated Event"
	err = s.UpdateEvent(event.ID, updatedEvent)
	require.NoError(t, err)

	// Проверяем, что событие обновлено
	events, err := s.GetEventsForDay(event.StartTime)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, updatedEvent, events[0])
}

// TestDeleteEvent проверяет удаление события.
func TestDeleteEvent(t *testing.T) {
	s := New()
	event := models.Event{
		ID:        uuid.New(),
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}

	// Создаем событие
	err := s.CreateEvent(event)
	require.NoError(t, err)

	// Удаляем событие
	err = s.DeleteEvent(event.ID)
	require.NoError(t, err)

	// Проверяем, что событие удалено
	events, err := s.GetEventsForDay(event.StartTime)
	require.NoError(t, err)
	require.Len(t, events, 0)
}

// TestGetEventsForDay проверяет получение событий на день.
func TestGetEventsForDay(t *testing.T) {
	s := New()
	now := time.Now()
	event1 := models.Event{
		ID:        uuid.New(),
		Title:     "Event 1",
		StartTime: now,
		EndTime:   now.Add(time.Hour),
	}
	event2 := models.Event{
		ID:        uuid.New(),
		Title:     "Event 2",
		StartTime: now.Add(24 * time.Hour), // Событие на следующий день
		EndTime:   now.Add(25 * time.Hour),
	}

	// Создаем события
	err := s.CreateEvent(event1)
	require.NoError(t, err)
	err = s.CreateEvent(event2)
	require.NoError(t, err)

	// Получаем события на сегодня
	events, err := s.GetEventsForDay(now)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, event1, events[0])
}

// TestGetEventsForWeek проверяет получение событий на неделю.
func TestGetEventsForWeek(t *testing.T) {
	s := New()
	now := time.Now()
	event1 := models.Event{
		ID:        uuid.New(),
		Title:     "Event 1",
		StartTime: now,
		EndTime:   now.Add(time.Hour),
	}
	event2 := models.Event{
		ID:        uuid.New(),
		Title:     "Event 2",
		StartTime: now.Add(7 * 24 * time.Hour), // Событие через неделю
		EndTime:   now.Add(7*24*time.Hour + time.Hour),
	}

	// Создаем события
	err := s.CreateEvent(event1)
	require.NoError(t, err)
	err = s.CreateEvent(event2)
	require.NoError(t, err)

	// Получаем события на текущую неделю
	events, err := s.GetEventsForWeek(now)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, event1, events[0])
}

// TestGetEventsForMonth проверяет получение событий на месяц.
func TestGetEventsForMonth(t *testing.T) {
	s := New()
	now := time.Now()
	event1 := models.Event{
		ID:        uuid.New(),
		Title:     "Event 1",
		StartTime: now,
		EndTime:   now.Add(time.Hour),
	}
	event2 := models.Event{
		ID:        uuid.New(),
		Title:     "Event 2",
		StartTime: now.Add(30 * 24 * time.Hour), // Событие через месяц
		EndTime:   now.Add(30*24*time.Hour + time.Hour),
	}

	// Создаем события
	err := s.CreateEvent(event1)
	require.NoError(t, err)
	err = s.CreateEvent(event2)
	require.NoError(t, err)

	// Получаем события на текущий месяц
	events, err := s.GetEventsForMonth(now)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, event1, events[0])
}

// TestConcurrency проверяет конкурентную безопасность хранилища.
func TestConcurrency(t *testing.T) {
	s := New()
	var wg sync.WaitGroup
	numGoroutines := 100

	// Создаем события в нескольких горутинах
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			event := models.Event{
				ID:        uuid.New(),
				Title:     "Event " + string(rune(i)),
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			}
			err := s.CreateEvent(event)
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()

	// Проверяем, что все события созданы
	events, err := s.GetEventsForMonth(time.Now())
	require.NoError(t, err)
	require.Len(t, events, numGoroutines)
}
