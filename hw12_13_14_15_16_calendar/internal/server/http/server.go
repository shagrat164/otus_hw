package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/server/handlers"
)

// Server отвечает за запуск HTTP-сервера.
type Server struct {
	host string
	port string
	logg Logger
	app  Application
	srv  *http.Server
}

// Logger - интерфейс логирования, используемый сервером.
type Logger interface {
	// Info выводит сообщение уровня INFO.
	Info(msg string)

	// Error выводит сообщение уровня ERROR.
	Error(msg string)
}

// Application - интерфейс приложения.
type Application interface { // TODO
}

// NewServer создаёт HTTP-сервер.
func NewServer(logger Logger, app Application, cfg *config.Config) *Server {
	addr := net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &Server{
		host: cfg.HTTP.Host,
		port: cfg.HTTP.Port,
		logg: logger,
		app:  app,
		srv: &http.Server{ //nolint:gosec
			Addr:    addr,
			Handler: nil,
		},
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	s.routes(mux)

	s.srv.Handler = loggingMiddleware(s.logg, mux)

	go func() {
		s.logg.Info("Starting HTTP server on " + s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil {
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

func (s *Server) routes(mux *http.ServeMux) {
	mux.HandleFunc("/status", handlers.StatusHandler())
}
