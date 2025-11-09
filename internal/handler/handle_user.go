package handle

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (h *UserHanlder) 	HandleCreateUser( w http.ResponseWriter, r *http.Request){
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "error_bad_request", "Format JSON tidak falid!")
		return
	}

	err = validate.Struct(req)

	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
            errors := formatValidationErrors(validationErrors)
            response.WriteJSON(w, http.StatusUnprocessableEntity, "error_validation", errors)
            return
        }
	}

	newUser, err := h.userService.Create(r.Context(), req.Name, req.Email, req.Password)

    if err != nil {
        if err.Error() == "email sudah terdaftar" {
            response.WriteJSON(w, http.StatusConflict, "error_conflict", err.Error())
        } else if err.Error() == "password minimal 8 karakter" {
            response.WriteJSON(w, http.StatusUnprocessableEntity, "error_validation", err.Error())
        } else {
            response.WriteJSON(w, http.StatusInternalServerError, "error_internal", err.Error())
        }
        return
    }
	response.WriteJSON(w, http.StatusCreated, "success_created", newUser)
}

func (h *UserHanlder) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")
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

func formatValidationErrors(errs validator.ValidationErrors) []ValidationErrorResponse {
    var errors []ValidationErrorResponse
    
    for _, err := range errs {
        field := err.Field()
        tag := err.Tag()
        
        message := fmt.Sprintf("Field '%s' gagal validasi pada tag '%s'", field, tag)
        switch tag {
        case "required":
            message = fmt.Sprintf("Field '%s' tidak boleh kosong", field)
        case "email":
            message = fmt.Sprintf("Field '%s' harus berupa format email yang valid", field)
        case "min":
            message = fmt.Sprintf("Field '%s' minimal harus %s karakter", field, err.Param())
        }

        errors = append(errors, ValidationErrorResponse{
            Field:   field,
            Message: message,
        })
    }
    return errors
}