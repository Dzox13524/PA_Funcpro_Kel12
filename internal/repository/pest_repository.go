package repository

import (
	"context"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PestRepository interface {
	Create(ctx context.Context, report *domain.PestReport) error
	GetAll(ctx context.Context) ([]domain.PestReport, error)
	GetByID(ctx context.Context, id string) (*domain.PestReport, error)
	IncrementVerification(ctx context.Context, id string) error
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

func (r *pestRepository) IncrementVerification(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.PestReport{}).
		Where("id = ?", id).
		UpdateColumn("verification_count", gorm.Expr("verification_count + ?", 1)).Error
}