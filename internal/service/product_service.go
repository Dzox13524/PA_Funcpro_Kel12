package service

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type CreateProductServiceFunc func(ctx context.Context, sellerID string, req domain.CreateProductRequest) (domain.Product, error)
type GetAllProductsServiceFunc func(ctx context.Context) ([]domain.Product, error)
type GetProductByIDServiceFunc func(ctx context.Context, id string) (domain.Product, error)
type UpdateProductServiceFunc func(ctx context.Context, id string, req domain.UpdateProductRequest, sellerID string) (domain.Product, error)
type DeleteProductServiceFunc func(ctx context.Context, id string, sellerID string) error
type UploadProductImageFunc func(ctx context.Context, id string, sellerID string, file io.Reader) (string, error)
type SearchProductsServiceFunc func(ctx context.Context, query, category, location string) ([]domain.Product, error)
type GetMetaCropsServiceFunc func(ctx context.Context) ([]string, error)
type GetMetaRegionsServiceFunc func(ctx context.Context) ([]string, error)

func NewCreateProductService(createRepo repository.CreateProductRepoFunc) CreateProductServiceFunc {
	return func(ctx context.Context, sellerID string, req domain.CreateProductRequest) (domain.Product, error) {
		newProduct := domain.Product{
			ID:          uuid.New().String(),
			SellerID:    sellerID, 
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
			Category:    req.Category,
			Location:    req.Location,
		}
		return createRepo(ctx, newProduct)
	}
}

func NewGetAllProductsService(getAllRepo repository.GetAllProductsRepoFunc) GetAllProductsServiceFunc {
	return func(ctx context.Context) ([]domain.Product, error) {
		return getAllRepo(ctx)
	}
}

func NewGetProductByIDService(getByIDRepo repository.GetProductByIDRepoFunc) GetProductByIDServiceFunc {
	return func(ctx context.Context, id string) (domain.Product, error) {
		return getByIDRepo(ctx, id)
	}
}

func NewUpdateProductService(updateRepo repository.UpdateProductRepoFunc, getByIDRepo repository.GetProductByIDRepoFunc) UpdateProductServiceFunc {
	return func(ctx context.Context, id string, req domain.UpdateProductRequest, sellerID string) (domain.Product, error) {
		existingProduct, err := getByIDRepo(ctx, id)
		if err != nil {
			return domain.Product{}, err
		}

		if existingProduct.SellerID != sellerID {
			return domain.Product{}, errors.New("anda bukan pemilik produk ini")
		}

		updates := map[string]interface{}{}
		if req.Name != "" { updates["name"] = req.Name }
		if req.Description != "" { updates["description"] = req.Description }
		if req.Price > 0 { updates["price"] = req.Price }
		if req.Stock >= 0 { updates["stock"] = req.Stock }
		updates["updated_at"] = time.Now()

		return updateRepo(ctx, id, updates)
	}
}

func NewDeleteProductService(deleteRepo repository.DeleteProductRepoFunc, getByIDRepo repository.GetProductByIDRepoFunc) DeleteProductServiceFunc {
	return func(ctx context.Context, id string, sellerID string) error {
		existingProduct, err := getByIDRepo(ctx, id)
		if err != nil {
			return err
		}
		if existingProduct.SellerID != sellerID {
			return errors.New("anda bukan pemilik produk ini")
		}
		
		return deleteRepo(ctx, id)
	}
}

func NewUploadProductImageService(getByIDRepo repository.GetProductByIDRepoFunc, updateRepo repository.UpdateProductRepoFunc) UploadProductImageFunc {
	return func(ctx context.Context, id string, sellerID string, file io.Reader) (string, error) {
		
		product, err := getByIDRepo(ctx, id)
		if err != nil { return "", err }
		if product.SellerID != sellerID {
			return "", errors.New("unauthorized: bukan barang anda")
		}
		
		cldName := "djf8xdry2"
		cldKey := "688323824432518"
		cldSecret := "WfBMEpA8JBB3xE835rHkKhBKl_A"
		cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)
		if err != nil {
			return "", err
		}
		uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
			Folder: "toko-belajar", 
			PublicID: id, 
		})
		if err != nil {
			return "", err
		}
		realURL := uploadResult.SecureURL
		_, err = updateRepo(ctx, id, map[string]interface{}{"image_url": realURL})
		
		return realURL, err
	}
}
func NewSearchProductsService(searchRepo repository.SearchProductsRepoFunc) SearchProductsServiceFunc {
	return func(ctx context.Context, query, category, location string) ([]domain.Product, error) {
		return searchRepo(ctx, query, category, location)
	}
}
func NewGetMetaCropsService(getMetaCropsRepo repository.GetMetaCropsRepoFunc) GetMetaCropsServiceFunc {
	return func(ctx context.Context) ([]string, error) {
		return getMetaCropsRepo(ctx)
	}
}
func NewGetMetaRegionsService(getMetaRegionsRepo repository.GetMetaRegionsRepoFunc) GetMetaRegionsServiceFunc {
	return func(ctx context.Context) ([]string, error) {
		return getMetaRegionsRepo(ctx)
	}
}