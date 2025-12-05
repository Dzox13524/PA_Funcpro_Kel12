package handle

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
}

func HandleCreateUser(createUser service.CreateUserFunc) func(*http.Request) (int, any, error) {
	return func(r *http.Request) (int, any, error) {
		var req CreateUserRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Format JSON tidak valid!")
		}

		newUser, err := createUser(r.Context(), req.Name, req.Email, req.Password)
		if err != nil {
			if err.Error() == "email sudah terdaftar" {
				return 0, nil, response.NewAPIError(http.StatusConflict, err.Error())
			}
			return 0, nil, err
		}

		return http.StatusCreated, newUser, nil
	}
}

func HandleGetUserByID(getUserByID service.GetUserByIDFunc) func(*http.Request) (int, any, error) {
	return func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")
		if id == "" {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "ID User diperlukan")
		}

		user, err := getUserByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "user tidak ditemukan" {
				return 0, nil, response.NewAPIError(http.StatusNotFound, "User tidak ditemukan")
			}
			return 0, nil, err
		}
		return http.StatusOK, user, nil
	}
}