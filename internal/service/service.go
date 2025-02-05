package service

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/spanwalla/docker-monitoring-backend/internal/repository"
	"github.com/spanwalla/docker-monitoring-backend/pkg/hasher"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

// PingerCreateInput -.
type PingerCreateInput struct {
	Name     string
	Password string
}

// PingerGenerateTokenInput -.
type PingerGenerateTokenInput struct {
	Name     string
	Password string
}

// Pinger -.
type Pinger interface {
	CreatePinger(ctx context.Context, input PingerCreateInput) (int, error)
	GenerateToken(ctx context.Context, input PingerGenerateTokenInput) (string, error)
	ParseToken(token string) (int, error)
}

// ReportStoreInput -.
type ReportStoreInput struct {
	PingerId int
	Report   []PingResult
}

// PingResult -.
type PingResult struct {
	ContainerId string    `json:"id" validate:"required,alphanum"`
	Ip          string    `json:"ip" validate:"required,ipv4"`
	Latency     int       `json:"latency_ms" validate:"required,number"`
	Timestamp   time.Time `json:"timestamp" validate:"required,datetime"`
}

// ReportOutput -.
type ReportOutput struct {
	Id         int          `json:"id"`
	PingerName string       `json:"pinger_name"`
	Content    pgtype.JSONB `json:"content"`
	CreatedAt  time.Time    `json:"created_at"`
}

// Report -.
type Report interface {
	Store(ctx context.Context, input ReportStoreInput) error
	GetActualReports(ctx context.Context) ([]ReportOutput, error)
}

// Services -.
type Services struct {
	Pinger
	Report
}

// Dependencies -.
type Dependencies struct {
	Repos    *repository.Repositories
	Hasher   hasher.PasswordHasher
	SignKey  string
	TokenTTL time.Duration
}

// NewServices -.
func NewServices(deps Dependencies) *Services {
	return &Services{
		Pinger: NewPingerService(deps.Repos.Pinger, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Report: NewReportService(deps.Repos.Report),
	}
}
