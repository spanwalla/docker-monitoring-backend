package service

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgtype"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/docker-monitoring-backend/internal/broker"
	"github.com/spanwalla/docker-monitoring-backend/internal/entity"
	"github.com/spanwalla/docker-monitoring-backend/internal/repository"
)

// ReportService -.
type ReportService struct {
	reportRepo repository.Report
	publisher  broker.Publisher
}

// NewReportService -.
func NewReportService(reportRepo repository.Report, publisher broker.Publisher) *ReportService {
	return &ReportService{reportRepo, publisher}
}

// PublishToQueue -.
func (s *ReportService) PublishToQueue(_ context.Context, input ReportStoreInput) error {
	// Convert []PingResult to json
	jsonBytes, err := json.Marshal(input.Report)
	if err != nil {
		log.Errorf("ReportService.PublishToQueue - json.Marshal: %v", err)
		return ErrCannotConvertJson
	}

	// Convert json to pgtype.JSONB
	var content pgtype.JSONB
	err = content.Set(jsonBytes)
	if err != nil {
		log.Errorf("ReportService.PublishToQueue - json.Set: %v", err)
		return ErrCannotConvertJson
	}

	report := entity.Report{
		PingerId: input.PingerId,
		Content:  content,
	}

	reportBytes, err := json.Marshal(report)
	if err != nil {
		log.Errorf("ReportService.PublishToQueue - json.Marshal: %v", err)
		return ErrCannotConvertJson
	}

	err = s.publisher.Publish(reportBytes)
	if err != nil {
		log.Errorf("ReportService.PublishToQueue - Publish: %v", err)
		return ErrCannotStoreReport
	}
	return nil
}

// Store -.
func (s *ReportService) Store(ctx context.Context, deliveryBody []byte) error {
	report := entity.Report{}
	if err := json.Unmarshal(deliveryBody, &report); err != nil {
		log.Errorf("ReportService.Store - json.Unmarshal: %v", err)
		return ErrCannotConvertJson
	}

	if err := s.reportRepo.CreateReport(ctx, report); err != nil {
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
