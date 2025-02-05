package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/docker-monitoring-backend/config"
	_ "github.com/spanwalla/docker-monitoring-backend/docs"
	v1 "github.com/spanwalla/docker-monitoring-backend/internal/controller/http/v1"
	"github.com/spanwalla/docker-monitoring-backend/internal/repository"
	"github.com/spanwalla/docker-monitoring-backend/internal/service"
	"github.com/spanwalla/docker-monitoring-backend/pkg/hasher"
	"github.com/spanwalla/docker-monitoring-backend/pkg/httpserver"
	"github.com/spanwalla/docker-monitoring-backend/pkg/postgres"
	"github.com/spanwalla/docker-monitoring-backend/pkg/validator"
	"os"
	"os/signal"
	"syscall"
)

// @title Docker Monitoring Service
// @version 1.0
// @description This is a service for storing and showing docker container's reports.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description JSON Web Token

// Run creates objects via constructors
func Run() {
	// Config
	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok || len(configPath) == 0 {
		log.Fatal("app - os.LookupEnv: CONFIG_PATH is empty")
	}

	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("app - config.New: %w", err))
	}

	// Logger
	setLogrus(cfg.Log.Level)
	log.Info("Config read")

	// Postgres
	log.Info("Connecting to postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Services and repos
	log.Info("Initializing services and repos...")
	services := service.NewServices(service.Dependencies{
		Repos:    repository.NewRepositories(pg),
		Hasher:   hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	})

	// Echo handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()
	v1.ConfigureRouter(handler, services)

	// RabbitMQ RPC Server
	// rmqRouter := amqrpc.NewRouter()
	// rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter)
	// if err != nil {
	//	log.Fatalf("app - Run - rmqServer - server.New: %v", err)
	// }

	// HTTP Server
	log.Info("Starting HTTP server...")
	log.Debugf("Server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Errorf("app - Run - httpServer.Notify: %v", err)
		// case err = <-rmqServer.Notify():
		//	log.Errorf("app - Run - rmqServer.Notify: %v", err)
	}

	// Graceful shutdown
	log.Info("Shutting down...")

	err = httpServer.Shutdown()
	if err != nil {
		log.Errorf("app - Run - httpServer.Shutdown: %v", err)
	}

	// err = rmqServer.Shutdown()
	// if err != nil {
	//	log.Errorf("app - Run - rmqServer.Shutdown: %v", err)
	// }
}
