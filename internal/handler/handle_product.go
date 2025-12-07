package handle

import (
	"encoding/json"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

// POST /api/v1/market/products
func HandleCreateProduct(createService service.CreateProductServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil User ID dari Context (diset oleh Middleware Auth)
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			response.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		// 2. Decode JSON body
		var req domain.CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}

		// 3. Panggil Service
		product, err := createService(r.Context(), userID, req)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusCreated, "product created", product)
	}
}

// GET /api/v1/market/products
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

// GET /api/v1/market/products/{id}
func HandleGetProductByID(getByIDService service.GetProductByIDServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id") // Fitur baru Go 1.22
		product, err := getByIDService(r.Context(), id)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, "product not found", nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "success", product)
	}
}

// PATCH /api/v1/market/products/{id}
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

// DELETE /api/v1/market/products/{id}
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

// POST /api/v1/market/products/{id}/upload-url
func HandleUploadProductImage(uploadService service.UploadProductImageFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil ID User & ID Produk
		userID := middleware.GetUserIDFromContext(r.Context())
		id := r.PathValue("id")

		// 2. Parse Multipart Form (Maksimal 10 MB)
		// Ini wajib agar backend bisa membaca file yang dikirim
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "File terlalu besar", nil)
			return
		}

		// 3. Ambil file dari form dengan key "image"
		// Di Postman nanti pilih Body -> form-data -> Key: "image", Type: File
		file, _, err := r.FormFile("image")
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Gagal membaca file gambar", nil)
			return
		}
		defer file.Close() // Tutup file setelah selesai

		// 4. Panggil Service Upload
		url, err := uploadService(r.Context(), id, userID, file)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		// 5. Return Sukses
		response.WriteJSON(w, http.StatusOK, "upload_success", map[string]string{"image_url": url})
	}
}