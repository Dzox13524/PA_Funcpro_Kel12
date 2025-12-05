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

func HandleCreateUser(createUser service.CreateUserFunc) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		var req CreateUserRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return response.NewAPIError(http.StatusBadRequest, "Format JSON tidak valid!")
		}

		newUser, err := createUser(r.Context(), req.Name, req.Email, req.Password)
		if err != nil {
			if err.Error() == "email sudah terdaftar" {
				return response.NewAPIError(http.StatusConflict, err.Error())
			}
			return err
		}

		response.WriteJSON(w, http.StatusCreated, "success_created", newUser)
		return nil
	}
}

func HandleGetUserByID(getUserByID service.GetUserByIDFunc) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := r.PathValue("id")
		if id == "" {
			return response.NewAPIError(http.StatusBadRequest, "ID User diperlukan")
		}

		user, err := getUserByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "user tidak ditemukan" {
				return response.NewAPIError(http.StatusNotFound, "User tidak ditemukan")
			}
			return err
		}
		response.WriteJSON(w, http.StatusOK, "success", user)
		return nil
	}
}