package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/Piccadilly98/subscription_service/docs"
	"github.com/Piccadilly98/subscription_service/internal/config"
	"github.com/Piccadilly98/subscription_service/internal/errors_checker"
	"github.com/Piccadilly98/subscription_service/internal/handlers"
	"github.com/Piccadilly98/subscription_service/internal/service"
	"github.com/Piccadilly98/subscription_service/internal/storage/cache"
	"github.com/Piccadilly98/subscription_service/internal/storage/data_base"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	router    chi.Router
	server    *http.Server
	service   *service.Service
	isStarted bool
	logger    *log.Logger
	ch        chan error
}

func InitServer() (*Server, error) {
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	db, err := data_base.NewDB(config.ConnectionStr)
	if err != nil {
		return nil, err
	}
	var c *cache.Cache
	if config.NeededCache {
		c = cache.NewCache(time.Second * time.Duration(config.CacheTTLInSecond))
	}

	serv, err := service.NewService(db, c)
	if err != nil {
		return nil, err
	}

	ew := errors_checker.NewErrorWorker(config.LoggingUserError)
	r := chi.NewRouter()

	getSummary := handlers.NewGetSummaryHandler(serv, ew)
	updateHandler := handlers.NewUpdateHandler(serv, ew)
	getHandler := handlers.NewGetHandler(serv, ew)
	createHandler := handlers.NewCreateHandler(serv, ew)
	deleteHandler := handlers.NewDeleteHandler(serv, ew)
	healthHandler := handlers.NewHealthHandler(serv, ew)

	r.Get("/health-check", healthHandler.Handler)
	r.Get("/subscriptions/summary", getSummary.Handler)
	r.Delete("/subscriptions/{id}", deleteHandler.Handler)
	r.Put("/subscriptions/{id}", updateHandler.Handler)
	r.Post("/subscriptions", createHandler.Handler)
	r.Get("/subscriptions/{id}", getHandler.Handler)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	httpServer := &http.Server{
		Addr:    config.ServerAddr + ":" + config.ServerPort,
		Handler: r,
	}
	return &Server{
		router:  r,
		server:  httpServer,
		service: serv,
		ch:      make(chan error),
		logger:  log.New(os.Stdout, "[SERVER INFO]", log.Ldate|log.Ltime),
	}, nil
}

func (s *Server) Start() (chan error, error) {
	if s.isStarted {
		return nil, fmt.Errorf("server is already running")
	}
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			s.logger.Printf("ERROR: Server failed: %v", err)
			s.ch <- err
			close(s.ch)
		}
	}()
	s.isStarted = true
	s.logger.Printf("INFO: server start in %s\n", s.server.Addr)
	return s.ch, nil
}
