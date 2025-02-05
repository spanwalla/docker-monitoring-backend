package service

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgtype"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/docker-monitoring-backend/internal/entity"
	"github.com/spanwalla/docker-monitoring-backend/internal/repository"
)

// ReportService -.
type ReportService struct {
	reportRepo repository.Report
}

// NewReportService -.
func NewReportService(reportRepo repository.Report) *ReportService {
	return &ReportService{reportRepo}
}

// Store -.
func (s *ReportService) Store(ctx context.Context, input ReportStoreInput) error {
	// Convert []PingResult to json
	jsonBytes, err := json.Marshal(input.Report)
	if err != nil {
		log.Errorf("ReportService.Store - json.Marshal: %v", err)
		return ErrCannotConvertJson
	}

	// Convert json to pgtype.JSONB
	var content pgtype.JSONB
	err = content.Set(jsonBytes)
	if err != nil {
		log.Errorf("ReportService.Store - json.Set: %v", err)
		return ErrCannotConvertJson
	}

	report := entity.Report{
		PingerId: input.PingerId,
		Content:  content,
	}

	err = s.reportRepo.CreateReport(ctx, report)
	if err != nil {
		log.Errorf("ReportService.Store - s.reportRepo.CreateReport: %v", err)
		return ErrCannotStoreReport
	}

	return nil
}

// GetActualReports -.
func (s *ReportService) GetActualReports(ctx context.Context) ([]ReportOutput, error) {
	reports, pingerNames, err := s.reportRepo.GetLatestReportByEveryPinger(ctx)
	if err != nil {
		log.Errorf("ReportService.GetActualReports - s.reportRepoGetLatestReportByEveryPinger: %v", err)
		return []ReportOutput{}, ErrCannotGetReports
	}

	output := make([]ReportOutput, 0, len(reports))
	for i, report := range reports {
		output = append(output, ReportOutput{
			Id:         report.Id,
			PingerName: pingerNames[i],
			Content:    report.Content,
			CreatedAt:  report.CreatedAt,
		})
	}

	return output, nil
}
