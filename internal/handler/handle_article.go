package handle

import (
	"encoding/json"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

func HandleCreateArticle(createService service.CreateArticleServiceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		var req domain.CreateArticleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}

		article, err := createService(r.Context(), userID, req)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusCreated, "article created", article)
	}
}

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