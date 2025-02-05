package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spanwalla/docker-monitoring-backend/internal/entity"
	"github.com/spanwalla/docker-monitoring-backend/pkg/postgres"
)

// PingerRepo -.
type PingerRepo struct {
	*postgres.Postgres
}

// NewPingerRepo -.
func NewPingerRepo(pg *postgres.Postgres) *PingerRepo {
	return &PingerRepo{pg}
}

// CreatePinger -.
func (p *PingerRepo) CreatePinger(ctx context.Context, pinger entity.Pinger) (int, error) {
	sql, args, _ := p.Builder.
		Insert("pingers").
		Columns("name", "password").
		Values(pinger.Name, pinger.Password).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := p.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("PingerRepo.CreatePinger - QueryRow: %w", err)
	}

	return id, nil
}

// GetPingerByNameAndPassword -.
func (p *PingerRepo) GetPingerByNameAndPassword(ctx context.Context, name, password string) (entity.Pinger, error) {
	sql, args, _ := p.Builder.
		Select("id, name, password, created_at").
		From("pingers").
		Where("name = ? AND password = ?", name, password).
		ToSql()

	var pinger entity.Pinger
	err := p.Pool.QueryRow(ctx, sql, args...).Scan(
		&pinger.Id,
		&pinger.Name,
		&pinger.Password,
		&pinger.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Pinger{}, ErrNotFound
		}
		return entity.Pinger{}, fmt.Errorf("PingerRepo.GetPingerByNameAndPassword - QueryRow: %w", err)
	}

	return pinger, nil
}

// GetPingerByName -.
func (p *PingerRepo) GetPingerByName(ctx context.Context, name string) (entity.Pinger, error) {
	return entity.Pinger{}, nil
}

// GetPingerById -.
func (p *PingerRepo) GetPingerById(ctx context.Context, id int) (entity.Pinger, error) {
	return entity.Pinger{}, nil
}
