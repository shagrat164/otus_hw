package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/models"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server отвечает за запуск HTTP-сервера.
type Server struct {
	httpHost string
	httpPort string
	grpcHost string
	grpcPort string
	logg     Logger
	app      Application
	srv      *http.Server
}

// Logger - интерфейс логирования, используемый сервером.
type Logger interface {
	// Info выводит сообщение уровня INFO.
	Info(msg string)

	// Error выводит сообщение уровня ERROR.
	Error(msg string)
}

// Application - интерфейс приложения.
type Application interface {
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

// NewServer создаёт HTTP-сервер.
func NewServer(logger Logger, app Application, cfg *config.Config) *Server {
	addr := net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &Server{
		httpHost: cfg.HTTP.Host,
		httpPort: cfg.HTTP.Port,
		grpcHost: cfg.GRPC.Host,
		grpcPort: cfg.GRPC.Port,
		logg:     logger,
		app:      app,
		srv: &http.Server{ //nolint:gosec
			Addr:    addr,
			Handler: nil,
		},
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	// Мультиплексор для gRPC-gateway
	gwMux := runtime.NewServeMux()

	// Адрес gRPC сервера из конфига
	grpcAddr := net.JoinHostPort(s.grpcHost, s.grpcPort)

	// Подключение к серверу gRPC
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.logg.Error(fmt.Sprintf("Failed to dial grpc server: %v", err.Error()))
		return err
	}
	defer conn.Close()

	err = pb.RegisterCalendarHandler(ctx, gwMux, conn)
	if err != nil {
		s.logg.Error(fmt.Sprintf("Failed to register gateway handler: %v", err.Error()))
		return err
	}

	s.srv.Handler = loggingMiddleware(s.logg, gwMux)

	go func() {
		s.logg.Info("Starting HTTP server on " + s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logg.Error("HTTP server error: " + err.Error())
		}
	}()

	<-ctx.Done() // Ожидание сигнала на завершение.
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	s.logg.Info("HTTP server stopped")
	return nil
}
