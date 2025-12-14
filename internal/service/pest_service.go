package service

import (
	"context"
	"time"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
)

type PestService interface {
	CreateReport(ctx context.Context, userID, pestName, desc, city, severity string) (*domain.PestReport, error)
	GetAllReports(ctx context.Context) ([]domain.PestReport, error)
	GetReportDetail(ctx context.Context, id string) (*domain.PestReport, error)
	VerifyReport(ctx context.Context, id string, userID string) error
}

type pestService struct {
	repo repository.PestRepository
}

func NewPestService(repo repository.PestRepository) PestService {
	return &pestService{repo: repo}
}

func (s *pestService) CreateReport(ctx context.Context, userID, pestName, desc, city, severity string) (*domain.PestReport, error) {
	newReport := &domain.PestReport{
		UserID:      userID,
		PestName:    pestName,
		Description: desc,
		City:        city,
		Severity:    severity,
		CreatedAt:   time.Now(),
		VerificationCount: 0,
	}
	if err := s.repo.Create(ctx, newReport); err != nil {
		return nil, err
	}
	return newReport, nil
}

func (s *pestService) GetAllReports(ctx context.Context) ([]domain.PestReport, error) {
	return s.repo.GetAll(ctx)
}

func (s *pestService) GetReportDetail(ctx context.Context, id string) (*domain.PestReport, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *pestService) VerifyReport(ctx context.Context, id string, userID string) error {
	return s.repo.VerifyReport(ctx, id, userID)
}