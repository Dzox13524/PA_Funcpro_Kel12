package handle

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CreateUserRequest struct{
	Name string
	Email string
	Password string
}

type ValidationErrorResponse struct {
	Field string
	Message string
}

var validate = validator.New()

type UserHanlder struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) *UserHanlder {
	return &UserHanlder{
		userService: userService,
	}
}

func (h *UserHanlder) HandleCreateUser(w http.ResponseWriter, r *http.Request) error { // Return error
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return response.NewAPIError(http.StatusBadRequest, "Format JSON tidak valid!")
	}
	newUser, err := h.userService.Create(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "email sudah terdaftar" {
			return response.NewAPIError(http.StatusConflict, err.Error())
		}
		return err
	}

	response.WriteJSON(w, http.StatusCreated, "success_created", newUser)
	return nil
}

func (h *UserHanlder) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.WriteJSON(w, http.StatusBadRequest, "error_bad_request", "ID User diperlukan")
        return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "user tidak ditemukan" {
            response.WriteJSON(w, http.StatusNotFound, "error_not_found", "User tidak ditemukan")
        } else {
            response.WriteJSON(w, http.StatusInternalServerError, "error_internal", err.Error())
        }
        return
    }
	response.WriteJSON(w, http.StatusOK, "success", user)
	}