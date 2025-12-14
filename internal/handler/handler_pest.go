package handle

import (
	"encoding/json"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

func HandleCreateAlert(svc service.PestService) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Unauthorized")
		}

		var req domain.CreatePestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid JSON")
		}

		res, err := svc.CreateReport(r.Context(), userID, req.PestName, req.Description, req.City, req.Severity)
		if err != nil {
			return 0, nil, err
		}

		return http.StatusCreated, res, nil
	})
}

func HandleGetMapData(svc service.PestService) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		res, err := svc.GetAllReports(r.Context())
		if err != nil {
			return 0, nil, err
		}
		return http.StatusOK, res, nil
	})
}

func HandleGetAlertDetail(svc service.PestService) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")
		res, err := svc.GetReportDetail(r.Context(), id)
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusNotFound, "Laporan tidak ditemukan")
		}
		return http.StatusOK, res, nil
	})
}

func HandleVerifyAlert(svc service.PestService) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")
		userID := middleware.GetUserIDFromContext(r.Context())
		
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Unauthorized")
		}
		
		if err := svc.VerifyReport(r.Context(), id, userID); err != nil {
			if err.Error() == "User sudah verifikasi laporan ini" {
				return 0, nil, response.NewAPIError(http.StatusConflict, "Anda sudah memvalidasi laporan ini")
			}
			return 0, nil, err
		}
		
		return http.StatusOK, nil, nil
	})
}