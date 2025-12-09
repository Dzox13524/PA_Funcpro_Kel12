package repository

import (
	"context"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateTransactionRepoFunc func(ctx context.Context, tx *domain.MarketTransaction) error
type GetTransactionByIDRepoFunc func(ctx context.Context, id string) (*domain.MarketTransaction, error)
type GetTransactionsByUserRepoFunc func(ctx context.Context, userID, transType string) ([]domain.MarketTransaction, error)
type UpdateTransactionStatusRepoFunc func(ctx context.Context, id, status string) error

func NewCreateTransactionRepository(db *gorm.DB) CreateTransactionRepoFunc {
	return func(ctx context.Context, tx *domain.MarketTransaction) error {
		if tx.ID == "" {
			tx.ID = uuid.New().String()
		}
		return db.WithContext(ctx).Create(tx).Error
	}
}

func NewGetTransactionByIDRepository(db *gorm.DB) GetTransactionByIDRepoFunc {
	return func(ctx context.Context, id string) (*domain.MarketTransaction, error) {
		var tx domain.MarketTransaction
		err := db.WithContext(ctx).
			Preload("Product").
			Preload("Buyer").
			First(&tx, "id = ?", id).Error
		return &tx, err
	}
}

func NewGetTransactionsByUserRepository(db *gorm.DB) GetTransactionsByUserRepoFunc {
	return func(ctx context.Context, userID, transType string) ([]domain.MarketTransaction, error) {
		var txs []domain.MarketTransaction
		err := db.WithContext(ctx).
			Preload("Product").
			Where("buyer_id = ? AND type = ?", userID, transType).
			Order("created_at desc").
			Find(&txs).Error
		return txs, err
	}
}

func NewUpdateTransactionStatusRepository(db *gorm.DB) UpdateTransactionStatusRepoFunc {
	return func(ctx context.Context, id, status string) error {
		return db.WithContext(ctx).Model(&domain.MarketTransaction{}).
			Where("id = ?", id).
			Update("status", status).Error
	}
}