package handle

import (
	"encoding/json"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

func HandleCreateQuestion(svc service.CreateQuestionFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Unauthorized")
		}

		var req domain.CreateQuestionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid JSON")
		}

		res, err := svc(r.Context(), userID, req.Title, req.Content, req.Category)
		if err != nil {
			return 0, nil, err
		}
		
		return http.StatusCreated, res, nil
	})
}

func HandleGetFeed(svc service.GetFeedFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		
		res, err := svc(r.Context(), userID)
		if err != nil {
			return 0, nil, err
		}
		
		return http.StatusOK, res, nil
	})
}

func HandleGetQuestionDetail(svc service.GetQuestionDetailFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")
		userID := middleware.GetUserIDFromContext(r.Context())
		
		res, err := svc(r.Context(), id, userID)
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusNotFound, "Question not found")
		}
		
		return http.StatusOK, res, nil
	})
}

func HandleAddAnswer(svc service.AddAnswerFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Login required")
		}
		
		id := r.PathValue("id")
		
		var req domain.CreateAnswerReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid JSON")
		}

		res, err := svc(r.Context(), userID, id, req.Content)
		if err != nil {
			return 0, nil, err
		}
		
		return http.StatusCreated, res, nil
	})
}

func HandleToggleLike(svc service.ToggleLikeFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Login required")
		}
		id := r.PathValue("id")
		
		isLiked, count, err := svc(r.Context(), userID, id)
		if err != nil {
			return 0, nil, err
		}
		
		return http.StatusOK, map[string]interface{}{
			"is_liked":    isLiked,
			"likes_count": count,
		}, nil
	})
}

func HandleToggleFav(svc service.ToggleFavFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Login required")
		}
		id := r.PathValue("id")
		
		isFav, err := svc(r.Context(), userID, id)
		if err != nil {
			return 0, nil, err
		}
		
		return http.StatusOK, map[string]interface{}{
			"is_favorited": isFav,
		}, nil
	})
}