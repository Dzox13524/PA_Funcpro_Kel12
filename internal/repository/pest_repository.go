package repository

import (
	"context"
	"fmt"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PestRepository interface {
	Create(ctx context.Context, report *domain.PestReport) error
	GetAll(ctx context.Context) ([]domain.PestReport, error)
	GetByID(ctx context.Context, id string) (*domain.PestReport, error)
	VerifyReport(ctx context.Context, id string ,userID string) error
}

type pestRepository struct {
	db *gorm.DB
}

func NewPestRepository(db *gorm.DB) PestRepository {
	return &pestRepository{db: db}
}

func (r *pestRepository) Create(ctx context.Context, report *domain.PestReport) error {
	if report.ID == "" {
		report.ID = uuid.New().String()
	}
	
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *pestRepository) GetAll(ctx context.Context) ([]domain.PestReport, error) {
	var reports []domain.PestReport
	err := r.db.WithContext(ctx).Order("created_at desc").Limit(100).Find(&reports).Error
	return reports, err
}

func (r *pestRepository) GetByID(ctx context.Context, id string) (*domain.PestReport, error) {
	var report domain.PestReport
	err := r.db.WithContext(ctx).First(&report, "id = ?", id).Error
	return &report, err
}

func (r *pestRepository) VerifyReport(ctx context.Context, reportID, userID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var verification domain.PestVerification
		err := tx.Where("user_id = ? AND pest_report_id = ?", userID, reportID).First(&verification).Error
		
		if err == nil {
			return fmt.Errorf("user already verified this report")
		}

		newVerification := domain.PestVerification{
			UserID:       userID,
			PestReportID: reportID,
		}
		
		if err := tx.Create(&newVerification).Error; err != nil {
			return err
		}
		if err := tx.Model(&domain.PestReport{}).
			Where("id = ?", reportID).
			UpdateColumn("verification_count", gorm.Expr("verification_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
}