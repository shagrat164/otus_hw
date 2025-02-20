package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	internalhttp "github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/server/http"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/storage/initstorage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Создание конфига.
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		// Если ошибка, конфиг по умолчанию.
		cfg = config.NewConfig()
	}

	// Создание логера.
	logg, _ := logger.New(cfg.Logger)

	// Создание хранилища для данных.
	storageData, err := initstorage.NewEventsStorage(cfg)
	if err != nil {
		// Завершение если ошибка.
		logg.Error(err.Error())
		return
	}

	// Создание приложения.
	calendar := app.New(logg, storageData, *cfg)

	// Создание HTTP-сервера.
	server := internalhttp.NewServer(logg, calendar, cfg)

	// Обработка сигналов для graceful shutdown.
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// Горутина для остановки HTTP-сервера.
	go func() {
		<-ctx.Done() // Ожидание сигнала на завершение.
		logg.Info("Shutting down HTTP server...")

		// Задержка для корректного завершения.
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), time.Second*3)
		defer cancelShutdown()

		if err := server.Stop(ctxShutdown); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	// Запуск сервера.
	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
