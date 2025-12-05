package middleware

import (
	"net/http"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
)

type APIFunc func(r *http.Request) (int, any, error)

func MakeHandler(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode, data, err := fn(r)

		if err != nil {
			if apiErr, ok := err.(response.APIError); ok {
				response.WriteJSON(w, apiErr.StatusCode, "error", apiErr.Message)
			} else {
				response.WriteJSON(w, http.StatusInternalServerError, "error_internal", err.Error())
			}
			return
		}
		response.WriteJSON(w, statusCode, "success", data)
	}
}