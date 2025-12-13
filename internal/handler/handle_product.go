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
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			response.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		var req domain.CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}
		product, err := createService(r.Context(), userID, req)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusCreated, "product created", product)
	}
}

func HandleGetAllProducts(getAllService service.GetAllProductsServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := getAllService(r.Context())
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "success", products)
	}
}

func HandleGetProductByID(getByIDService service.GetProductByIDServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		product, err := getByIDService(r.Context(), id)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, "product not found", nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "success", product)
	}
}

func HandleUpdateProduct(updateService service.UpdateProductServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		id := r.PathValue("id")

		var req domain.UpdateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}

		product, err := updateService(r.Context(), id, req, userID)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "product updated", product)
	}
}

func HandleDeleteProduct(deleteService service.DeleteProductServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		id := r.PathValue("id")

		if err := deleteService(r.Context(), id, userID); err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "product deleted", nil)
	}
}

func HandleUploadProductImage(uploadService service.UploadProductImageFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		id := r.PathValue("id")
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "File terlalu besar", nil)
			return
		}
		file, _, err := r.FormFile("image")
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Gagal membaca file gambar", nil)
			return
		}
		defer file.Close()
		url, err := uploadService(r.Context(), id, userID, file)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "upload_success", map[string]string{"image_url": url})
	}
}
func HandleSearchProducts(searchService service.SearchProductsServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")      
		category := r.URL.Query().Get("category") 
		location := r.URL.Query().Get("location")
		products, err := searchService(r.Context(), query, category, location)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "success", products)
	}
}
func HandleGetMetaCrops(getMetaCropsService service.GetMetaCropsServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		crops, err := getMetaCropsService(r.Context())	
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}	
		response.WriteJSON(w, http.StatusOK, "success", crops)
	}
}

func HandleGetMetaRegions(getMetaRegionsService service.GetMetaRegionsServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regions, err := getMetaRegionsService(r.Context())
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}	
		response.WriteJSON(w, http.StatusOK, "success", regions)
	}
}
