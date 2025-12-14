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
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		
		var req domain.CreateArticleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid JSON")
		}

		article, err := createService(r.Context(), userID, req)
		if err != nil {
			return 0, nil, err
		}

		return http.StatusCreated, article, nil
	})
}

func HandleGetAllArticles(getAllService service.GetAllArticlesServiceFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		articles, err := getAllService(r.Context())
		if err != nil {
			return 0, nil, err
		}
		
		return http.StatusOK, articles, nil
	})
}

func HandleGetArticleByID(getByIDService service.GetArticleByIDServiceFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")
		article, err := getByIDService(r.Context(), id)
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusNotFound, "article not found")
		}
		
		return http.StatusOK, article, nil
	})
}