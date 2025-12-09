package handle

import (
	"encoding/json"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

type CreateQuestionReq struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type CreateAnswerReq struct {
	Content string `json:"content"`
}

func HandleCreateQuestion(svc service.CreateQuestionFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			response.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		var req CreateQuestionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}

		res, err := svc(r.Context(), userID, req.Title, req.Content, req.Category)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusCreated, "Question created", res)
	}
}

func HandleGetFeed(svc service.GetFeedFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context()) // Boleh kosong jika guest
		res, err := svc(r.Context(), userID)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "Success", res)
	}
}

func HandleGetQuestionDetail(svc service.GetQuestionDetailFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		userID := middleware.GetUserIDFromContext(r.Context())
		
		res, err := svc(r.Context(), id, userID)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, "Question not found", nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "Success", res)
	}
}

func HandleAddAnswer(svc service.AddAnswerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			response.WriteJSON(w, http.StatusUnauthorized, "Login required", nil)
			return
		}
		
		id := r.PathValue("id")
		var req CreateAnswerReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}

		res, err := svc(r.Context(), userID, id, req.Content)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusCreated, "Answer added", res)
	}
}

func HandleToggleLike(svc service.ToggleLikeFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			response.WriteJSON(w, http.StatusUnauthorized, "Login required", nil)
			return
		}
		id := r.PathValue("id")
		
		isLiked, count, err := svc(r.Context(), userID, id)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		
		response.WriteJSON(w, http.StatusOK, "Success", map[string]interface{}{
			"is_liked": isLiked,
			"likes_count": count,
		})
	}
}

func HandleToggleFav(svc service.ToggleFavFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			response.WriteJSON(w, http.StatusUnauthorized, "Login required", nil)
			return
		}
		id := r.PathValue("id")
		
		isFav, err := svc(r.Context(), userID, id)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		
		response.WriteJSON(w, http.StatusOK, "Success", map[string]interface{}{
			"is_favorited": isFav,
		})
	}
}