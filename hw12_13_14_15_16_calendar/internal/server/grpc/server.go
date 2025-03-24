package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/models"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/server/grpc/pb"
	"github.com/shagrat164/otus_hw/hw12_13_14_15_16_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server отвечает за запуск gRPC-сервера.
type Server struct {
	host string
	port string
	logg Logger
	app  Application
	srv  *grpc.Server
}

type CalendarService struct {
	pb.UnimplementedCalendarServer
	app  Application
	logg Logger
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

func NewServer(logger Logger, app Application, cfg *config.Config) *Server {
	chainInterceptor := grpc.ChainUnaryInterceptor(
		unaryInterceptor(logger),
	)
	server := grpc.NewServer(chainInterceptor)

	service := &CalendarService{
		app:  app,
		logg: logger,
	}
	pb.RegisterCalendarServer(server, service)

	// Включаем reflection для удобства отладки с grpcurl
	reflection.Register(server)

	srv := &Server{
		host: cfg.GRPC.Host,
		port: cfg.GRPC.Port,
		logg: logger,
		app:  app,
		srv:  server,
	}

	return srv
}

// Start - запускает gRPC-сервер.
func (s *Server) Start() error {
	addr := net.JoinHostPort(s.host, s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logg.Info(fmt.Sprintf("gRPC server started on %s", addr))

	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}

// CreateEvent создаёт новое событие.
func (s *CalendarService) CreateEvent(_ context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	// Проверка на корректный Event
	reqEvent := req.GetEvent()
	if reqEvent == nil {
		return nil, status.Error(codes.InvalidArgument, "event is not specified")
	}

	s.logg.Info(fmt.Sprintf("%v\n", reqEvent))

	// Создать и заполнить данными новый Event
	event := models.NewEvent()
	event.Title = reqEvent.GetTitle()
	event.StartTime = reqEvent.GetStartTime().AsTime()
	event.EndTime = reqEvent.GetEndTime().AsTime()
	event.Description = reqEvent.GetDescription()
	event.UserID, _ = uuid.Parse(reqEvent.GetUserID().GetVal())
	event.Reminder = reqEvent.GetReminder().AsTime()

	// Создать Event
	err := s.app.CreateEvent(*event)
	if err != nil {
		if errors.Is(err, storage.ErrEventAlreadyExist) {
			return nil, status.Errorf(codes.NotFound, "event %s is exists", event.ID.String())
		}
		s.logg.Error(fmt.Sprintf("Fail create event %v", err))
		return nil, status.Errorf(codes.DataLoss, "event create err: %v", err)
	}

	eventID := new(pb.UUID)
	eventID.Val = event.ID.String()

	resp := &pb.CreateEventResponse{}
	resp.Id = eventID

	return resp, nil
}

// UpdateEvent обновляет существующее событие по id.
func (s *CalendarService) UpdateEvent(_ context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	// Проверка на корректный Id
	reqID := req.GetId()
	if reqID == nil {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}

	// Проверка на корректный Event
	reqEvent := req.GetEvent()
	if reqEvent == nil {
		return nil, status.Error(codes.InvalidArgument, "event is not specified")
	}

	// Полученый id перевести в uuid
	eventID, _ := uuid.Parse(reqID.GetVal())
	eventUserID, _ := uuid.Parse(reqEvent.GetUserID().GetVal())

	// Создать и заполнить данными новый Event
	event := models.NewEvent()
	event.ID = eventID
	event.Title = reqEvent.GetTitle()
	event.StartTime = reqEvent.GetStartTime().AsTime()
	event.EndTime = reqEvent.GetEndTime().AsTime()
	event.Description = reqEvent.GetDescription()
	event.UserID = eventUserID
	event.Reminder = reqEvent.GetReminder().AsTime()

	// Обновить Event
	err := s.app.UpdateEvent(eventID, *event)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Errorf(codes.NotFound, "event %s is not exists", reqID.GetVal())
		}

		return nil, status.Errorf(codes.Internal, "update events error: %v", err)
	}

	evID := new(pb.UUID)
	evID.Val = event.ID.String()

	resp := &pb.UpdateEventResponse{}
	resp.Id = evID

	return resp, nil
}

// DeleteEvent удаляет событие по id.
func (s *CalendarService) DeleteEvent(_ context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	// Проверка на корректный Id
	reqID := req.GetId()
	if reqID == nil {
		return nil, status.Error(codes.InvalidArgument, "id is not specified")
	}

	// Полученый id перевести в uuid
	eventID, _ := uuid.Parse(reqID.GetVal())

	err := s.app.DeleteEvent(eventID)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Errorf(codes.NotFound, "event %s is not exists", reqID.GetVal())
		}

		return nil, status.Errorf(codes.Internal, "delete events error: %v", err)
	}

	return &pb.DeleteEventResponse{}, nil
}

// GetEventsForDay возвращает список событий на указанный день.
func (s *CalendarService) GetEventsForDay( //nolint:dupl
	_ context.Context,
	req *pb.GetEventsForRequest,
) (*pb.GetEventsForResponse, error) {
	reqData := req.GetDate()
	if reqData == nil {
		return nil, status.Error(codes.InvalidArgument, "data is not specified")
	}

	events, err := s.app.GetEventsForDay(reqData.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get events day error: %v", err)
	}

	resp := make([]*pb.Event, 0, len(events))

	for _, event := range events {
		data := &pb.Event{
			Id:          &pb.UUID{Val: event.ID.String()},
			Title:       event.Title,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     timestamppb.New(event.EndTime),
			Description: event.Description,
			UserID:      &pb.UUID{Val: event.UserID.String()},
			Reminder:    timestamppb.New(event.Reminder),
		}
		resp = append(resp, data)
	}

	return &pb.GetEventsForResponse{Events: resp}, nil
}

// GetEventsForWeek возвращает список событий на неделю, начиная с указанной даты.
func (s *CalendarService) GetEventsForWeek( //nolint:dupl
	_ context.Context,
	req *pb.GetEventsForRequest,
) (*pb.GetEventsForResponse, error) {
	reqData := req.GetDate()
	if reqData == nil {
		return nil, status.Error(codes.InvalidArgument, "data is not specified")
	}

	events, err := s.app.GetEventsForWeek(reqData.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get events week error: %v", err)
	}

	resp := make([]*pb.Event, 0, len(events))

	for _, event := range events {
		data := &pb.Event{
			Id:          &pb.UUID{Val: event.ID.String()},
			Title:       event.Title,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     timestamppb.New(event.EndTime),
			Description: event.Description,
			UserID:      &pb.UUID{Val: event.UserID.String()},
			Reminder:    timestamppb.New(event.Reminder),
		}
		resp = append(resp, data)
	}

	return &pb.GetEventsForResponse{Events: resp}, nil
}

// GetEventsForMonth возвращает список событий на месяц, начиная с указанной даты.
func (s *CalendarService) GetEventsForMonth( //nolint:dupl
	_ context.Context,
	req *pb.GetEventsForRequest,
) (*pb.GetEventsForResponse, error) {
	reqData := req.GetDate()
	if reqData == nil {
		return nil, status.Error(codes.InvalidArgument, "data is not specified")
	}

	events, err := s.app.GetEventsForMonth(reqData.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get events month error: %v", err)
	}

	resp := make([]*pb.Event, 0, len(events))

	for _, event := range events {
		data := &pb.Event{
			Id:          &pb.UUID{Val: event.ID.String()},
			Title:       event.Title,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     timestamppb.New(event.EndTime),
			Description: event.Description,
			UserID:      &pb.UUID{Val: event.UserID.String()},
			Reminder:    timestamppb.New(event.Reminder),
		}
		resp = append(resp, data)
	}

	return &pb.GetEventsForResponse{Events: resp}, nil
}
