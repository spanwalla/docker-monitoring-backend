package repository

import (
	"context"
	"github.com/spanwalla/docker-monitoring-backend/internal/entity"
	"github.com/spanwalla/docker-monitoring-backend/pkg/postgres"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

// Pinger -.
type Pinger interface {
	CreatePinger(ctx context.Context, pinger entity.Pinger) (int, error)
	GetPingerByNameAndPassword(ctx context.Context, name, password string) (entity.Pinger, error)
	GetPingerByName(ctx context.Context, name string) (entity.Pinger, error)
	GetPingerById(ctx context.Context, id int) (entity.Pinger, error)
}

// Report -.
type Report interface {
	CreateReport(ctx context.Context, report entity.Report) error
	GetLatestReportByEveryPinger(ctx context.Context) ([]entity.Report, []string, error)
}

// Repositories -.
type Repositories struct {
	Pinger
	Report
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Pinger: NewPingerRepo(pg),
		Report: NewReportRepo(pg),
	}
}
