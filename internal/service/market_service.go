package service

import (
	"context"
	"errors"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
)

type CreateReservationFunc func(ctx context.Context, userID string, req domain.CreateOrderRequest) (*domain.MarketTransaction, error)
type CreateOrderFunc func(ctx context.Context, userID string, req domain.CreateOrderRequest) (*domain.MarketTransaction, error)
type GetTransactionDetailFunc func(ctx context.Context, id string) (*domain.MarketTransaction, error)
type GetUserTransactionsFunc func(ctx context.Context, userID string) ([]domain.MarketTransaction, error)
type UpdateTransactionStatusFunc func(ctx context.Context, id, newStatus string) error

func NewCreateReservationService(
	createTxRepo repository.CreateTransactionRepoFunc,
	getProdRepo repository.GetProductByIDRepoFunc,
	updateProdRepo repository.UpdateProductRepoFunc,
) CreateReservationFunc {
	return func(ctx context.Context, userID string, req domain.CreateOrderRequest) (*domain.MarketTransaction, error) {
		product, err := getProdRepo(ctx, req.ProductID)
		if err != nil {
			return nil, errors.New("produk tidak ditemukan")
		}

		if product.Stock < req.Quantity {
			return nil, errors.New("stok produk tidak mencukupi")
		}

		totalPrice := float64(req.Quantity) * product.Price

		tx := &domain.MarketTransaction{
			BuyerID:    userID,
			ProductID:  req.ProductID,
			Quantity:   req.Quantity,
			TotalPrice: totalPrice,
			Status:     domain.StatusPending,
			Type:       domain.TypeReservation,
			Note:       req.Note,
		}

		if err := createTxRepo(ctx, tx); err != nil {
			return nil, err
		}

		newStock := product.Stock - req.Quantity
		_, err = updateProdRepo(ctx, product.ID, map[string]interface{}{
			"stock": newStock,
		})

		if err != nil {
			return nil, errors.New("gagal update stok produk")
		}

		return tx, nil
	}
}

func NewCreateOrderService(
	createTxRepo repository.CreateTransactionRepoFunc,
	getProdRepo repository.GetProductByIDRepoFunc,
	updateProdRepo repository.UpdateProductRepoFunc,
) CreateOrderFunc {
	return func(ctx context.Context, userID string, req domain.CreateOrderRequest) (*domain.MarketTransaction, error) {
		product, err := getProdRepo(ctx, req.ProductID)
		if err != nil {
			return nil, errors.New("produk tidak ditemukan")
		}

		if product.Stock < req.Quantity {
			return nil, errors.New("stok produk tidak mencukupi")
		}

		totalPrice := float64(req.Quantity) * product.Price

		tx := &domain.MarketTransaction{
			BuyerID:    userID,
			ProductID:  req.ProductID,
			Quantity:   req.Quantity,
			TotalPrice: totalPrice,
			Status:     domain.StatusPending,
			Type:       domain.TypeOrder,
			Note:       req.Note,
		}

		if err := createTxRepo(ctx, tx); err != nil {
			return nil, err
		}

		newStock := product.Stock - req.Quantity
		_, err = updateProdRepo(ctx, product.ID, map[string]interface{}{
			"stock": newStock,
		})

		if err != nil {
			return nil, errors.New("gagal update stok produk")
		}

		return tx, nil
	}
}

func NewGetTransactionDetailService(getTxRepo repository.GetTransactionByIDRepoFunc) GetTransactionDetailFunc {
	return func(ctx context.Context, id string) (*domain.MarketTransaction, error) {
		return getTxRepo(ctx, id)
	}
}

func NewGetUserTransactionsService(getTxsRepo repository.GetTransactionsByUserRepoFunc, transType string) GetUserTransactionsFunc {
	return func(ctx context.Context, userID string) ([]domain.MarketTransaction, error) {
		return getTxsRepo(ctx, userID, transType)
	}
}

func NewUpdateTransactionStatusService(updateStatusRepo repository.UpdateTransactionStatusRepoFunc) UpdateTransactionStatusFunc {
	return func(ctx context.Context, id, newStatus string) error {
		return updateStatusRepo(ctx, id, newStatus)
	}
}