package repository

import (
	"context"
	"fmt"
	"github.com/spanwalla/docker-monitoring-backend/internal/entity"
	"github.com/spanwalla/docker-monitoring-backend/pkg/postgres"
)

// ReportRepo -.
type ReportRepo struct {
	*postgres.Postgres
}

// NewReportRepo -.
func NewReportRepo(pg *postgres.Postgres) *ReportRepo {
	return &ReportRepo{pg}
}

// CreateReport -.
func (r *ReportRepo) CreateReport(ctx context.Context, report entity.Report) error {
	sql, args, _ := r.Builder.
		Insert("reports").
		Columns("pinger_id", "content").
		Values(report.PingerId, report.Content).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ReportRepo.CreateReport - Exec: %w", err)
	}

	return nil
}

// GetLatestReportByEveryPinger -.
func (r *ReportRepo) GetLatestReportByEveryPinger(ctx context.Context) ([]entity.Report, []string, error) {
	sql, args, _ := r.Builder.
		Select("DISTINCT ON (reports.pinger_id) reports.id",
			"reports.pinger_id",
			"reports.content",
			"reports.created_at",
			"pingers.name AS pinger_name").
		From("reports").
		InnerJoin("pingers ON reports.pinger_id = pingers.id").
		OrderBy("reports.pinger_id", "reports.created_at DESC").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("ReportRepo.GetLatestReportByEveryPinger - Query: %w", err)
	}
	defer rows.Close()

	var reports []entity.Report
	var pingerNames []string
	for rows.Next() {
		var report entity.Report
		var pingerName string
		err = rows.Scan(&report.Id, &report.PingerId, &report.Content, &report.CreatedAt, &pingerName)
		if err != nil {
			return nil, nil, fmt.Errorf("ReportRepo.GetLatestReportByEveryPinger - rows.Scan: %w", err)
		}
		reports = append(reports, report)
		pingerNames = append(pingerNames, pingerName)
	}

	return reports, pingerNames, nil
}
