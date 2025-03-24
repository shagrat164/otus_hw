package initstorage

import (
	"context"
	"fmt"

	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	memorystorage "github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/storage/memory"
	sqlstorage "github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/storage/sql"
)

func getDSN(cfg config.DatabaseConf) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Pwd, cfg.Host, cfg.Port, cfg.DBName)
}

func NewEventsStorage(cfg *config.Config) (app.Storage, error) {
	switch cfg.Storage.Type {
	case "in-memory":
		return memorystorage.New(), nil
	case "sql":
		sqlStorage := sqlstorage.New(getDSN(cfg.Database))
		if err := sqlStorage.Connect(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to connect to storage: %w", err)
		}
		return sqlStorage, nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Storage.Type)
	}
}
