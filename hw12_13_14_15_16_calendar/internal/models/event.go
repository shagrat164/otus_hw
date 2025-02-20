package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID // Уникальный идентификатор события
	Title       string    // Заголовок
	StartTime   time.Time // Дата и время начала события
	EndTime     time.Time // Дата и время окончания события
	Description string    // Описание события, опционально
	UserID      uuid.UUID // ID пользователя, владельца события
	Reminder    time.Time // За сколько времени высылать уведомление, опционально
}
