package handle

import (
	"encoding/json"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

// Baris 15: Handler POST (Buat Artikel)
// Endpoint: POST /api/v1/admin/articles
func HandleCreateArticle(createService service.CreateArticleServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Cek User ID (karena rute ini diproteksi AuthMiddleware)
		userID := middleware.GetUserIDFromContext(r.Context())
		
		// 2. Decode JSON body
		var req domain.CreateArticleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}

		// 3. Panggil Service
		article, err := createService(r.Context(), userID, req)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusCreated, "article created", article)
	}
}

// Baris 40: Handler GET All (Lihat Daftar)
// Endpoint: GET /api/v1/articles
func HandleGetAllArticles(getAllService service.GetAllArticlesServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := getAllService(r.Context())
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "success", articles)
	}
}

// Baris 53: Handler GET By ID (Baca Detail)
// Endpoint: GET /api/v1/articles/{id}
func HandleGetArticleByID(getByIDService service.GetArticleByIDServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		article, err := getByIDService(r.Context(), id)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, "article not found", nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "success", article)
	}
}