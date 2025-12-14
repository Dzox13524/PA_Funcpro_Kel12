package handle

import (
	"encoding/json"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

func HandleCreateProduct(createService service.CreateProductServiceFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Unauthorized")
		}

		var req domain.CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid JSON")
		}
		
		product, err := createService(r.Context(), userID, req)
		if err != nil {
			return 0, nil, err
		}

		return http.StatusCreated, product, nil
	})
}

func HandleGetAllProducts(getAllService service.GetAllProductsServiceFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		products, err := getAllService(r.Context())
		if err != nil {
			return 0, nil, err
		}
		return http.StatusOK, products, nil
	})
}

func HandleGetProductByID(getByIDService service.GetProductByIDServiceFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")
		product, err := getByIDService(r.Context(), id)
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusNotFound, "product not found")
		}
		return http.StatusOK, product, nil
	})
}

func HandleUpdateProduct(updateService service.UpdateProductServiceFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		id := r.PathValue("id")

		var req domain.UpdateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid JSON")
		}

		product, err := updateService(r.Context(), id, req, userID)
		if err != nil {
			return 0, nil, err
		}
		
		return http.StatusOK, product, nil
	})
}

func HandleDeleteProduct(deleteService service.DeleteProductServiceFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		id := r.PathValue("id")

		if err := deleteService(r.Context(), id, userID); err != nil {
			return 0, nil, err
		}
		
		return http.StatusOK, nil, nil
	})
}

func HandleUploadProductImage(uploadService service.UploadProductImageFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		id := r.PathValue("id")
		
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "File terlalu besar")
		}
		
		file, _, err := r.FormFile("image")
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Gagal membaca file gambar")
		}
		defer file.Close()
		
		url, err := uploadService(r.Context(), id, userID, file)
		if err != nil {
			return 0, nil, err
		}
		
		return http.StatusOK, map[string]string{"image_url": url}, nil
	})
}

func HandleSearchProducts(searchService service.SearchProductsServiceFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		query := r.URL.Query().Get("q")      
		category := r.URL.Query().Get("category") 
		location := r.URL.Query().Get("location")
		
		products, err := searchService(r.Context(), query, category, location)
		if err != nil {
			return 0, nil, err
		}
		
		return http.StatusOK, products, nil
	})
}

func HandleGetMetaCrops(getMetaCropsService service.GetMetaCropsServiceFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		crops, err := getMetaCropsService(r.Context())    
		if err != nil {
			return 0, nil, err
		}   
		return http.StatusOK, crops, nil
	})
}

func HandleGetMetaRegions(getMetaRegionsService service.GetMetaRegionsServiceFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		regions, err := getMetaRegionsService(r.Context())
		if err != nil {
			return 0, nil, err
		}   
		return http.StatusOK, regions, nil
	})
}