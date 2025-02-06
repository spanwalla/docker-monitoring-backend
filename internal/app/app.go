package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/docker-monitoring-backend/config"
	_ "github.com/spanwalla/docker-monitoring-backend/docs"
	"github.com/spanwalla/docker-monitoring-backend/internal/broker"
	v1 "github.com/spanwalla/docker-monitoring-backend/internal/controller/http/v1"
	"github.com/spanwalla/docker-monitoring-backend/internal/repository"
	"github.com/spanwalla/docker-monitoring-backend/internal/service"
	"github.com/spanwalla/docker-monitoring-backend/pkg/hasher"
	"github.com/spanwalla/docker-monitoring-backend/pkg/httpserver"
	"github.com/spanwalla/docker-monitoring-backend/pkg/postgres"
	"github.com/spanwalla/docker-monitoring-backend/pkg/rabbitmq"
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

	// RabbitMQ
	rmq, err := rabbitmq.New(cfg.RMQ.URL)
	if err != nil {
		log.Fatalf("app - Run - rabbitmq.New: %v", err)
	}
	defer func() {
		if err = rmq.Close(); err != nil {
			log.Fatalf("app - Run - rmq.Close: %v", err)
		}
	}()

	pubChannel, err := rmq.Channel()
	if err != nil {
		log.Fatalf("app - Run - rmq.Channel (publisher): %v", err)
	}
	defer func() {
		if err = pubChannel.Close(); err != nil {
			log.Fatalf("app - Run - pubChannel.Close: %v", err)
		}
	}()

	pub, err := broker.NewRabbitMQPublisher(pubChannel, cfg.RMQ.ReportQueue)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - broker.NewRabbitMQPublisher: %w", err))
	}

	// Services and repos
	log.Info("Initializing services and repos...")
	services := service.NewServices(service.Dependencies{
		Repos:     repository.NewRepositories(pg),
		Hasher:    hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:   cfg.JWT.SignKey,
		TokenTTL:  cfg.JWT.TokenTTL,
		Publisher: pub,
	})

	// Echo handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()
	v1.ConfigureRouter(handler, services)

	// HTTP Server
	log.Info("Starting HTTP server...")
	log.Debugf("Server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Report consumer
	consChannel, err := rmq.Channel()
	if err != nil {
		log.Fatalf("app - Run - rmq.Channel (consumer): %v", err)
	}
	defer func(consChannel *amqp.Channel) {
		err = consChannel.Close()
		if err != nil {
			log.Fatalf("app - Run - consChannel.Close: %v", err)
		}
	}(consChannel)

	reportConsumer, err := broker.NewRabbitMQConsumer(consChannel, cfg.RMQ.ReportQueue, cfg.RMQ.ReportQueue, services.Store)
	if err != nil {
		log.Fatalf("app - Run - broker.NewRabbitMQConsumer: %v", err)
	}
	consumerErrChan := make(chan error, 1)
	go func() {
		if err = reportConsumer.Start(); err != nil {
			consumerErrChan <- fmt.Errorf("app - Run - reportConsumer.Start: %w", err)
		}
	}()

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Errorf("app - Run - httpServer.Notify: %v", err)
	case err = <-consumerErrChan:
		log.Errorf("app - Run - RabbitMQConsumer: %v", err)
	}

	// Graceful shutdown
	log.Info("Shutting down...")

	err = httpServer.Shutdown()
	if err != nil {
		log.Errorf("app - Run - httpServer.Shutdown: %v", err)
	}

	err = reportConsumer.Stop()
	if err != nil {
		log.Errorf("app - Run - reportConsumer.Stop: %v", err)
	}
}
