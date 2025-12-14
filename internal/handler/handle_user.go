package handle

import (
	"encoding/json"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

func HandleCreateUser(createUser service.CreateUserFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		var req domain.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Format JSON tidak valid!")
		}

		newUser, err := createUser(r.Context(), req.Name, req.Email, req.Password)
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusConflict, err.Error())
		}

		return http.StatusCreated, newUser, nil
	})
}

func HandleLogin(loginService service.LoginFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		var req domain.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid Request")
		}

		authResp, err := loginService(r.Context(), req.Email, req.Password)
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, err.Error())
		}

		return http.StatusOK, authResp, nil
	})
}

func HandleGetMe(getUserByID service.GetUserByIDFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Unauthorized")
		}

		user, err := getUserByID(r.Context(), userID)
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusNotFound, "User not found")
		}

		return http.StatusOK, user, nil
	})
}

func HandleUpdateMe(updateUser service.UpdateUserFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Unauthorized")
		}

		var req domain.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid JSON")
		}

		updatedUser, err := updateUser(r.Context(), userID, req.Name)
		if err != nil {
			return 0, nil, err
		}

		return http.StatusOK, updatedUser, nil
	})
}

func HandleGetUserByID(getUserByID service.GetUserByIDFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")
		user, err := getUserByID(r.Context(), id)
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusNotFound, "User tidak ditemukan")
		}
		return http.StatusOK, user, nil
	})
}