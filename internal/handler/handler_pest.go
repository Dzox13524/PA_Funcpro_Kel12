package handle

import (
	"encoding/json"
	"net/http"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

type CreatePestRequest struct {
	PestName    string `json:"pest_name"`
	Description string `json:"description"`
	City        string `json:"city"`
	Severity    string `json:"severity"`
}

func HandleCreateAlert(svc service.PestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		var req CreatePestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}
		res, err := svc.CreateReport(r.Context(), userID, req.PestName, req.Description, req.City, req.Severity)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusCreated, "Laporan Berhasil Dibuat", res)
	}
}

func HandleGetMapData(svc service.PestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := svc.GetAllReports(r.Context())
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "Data Peta Hama", res)
	}
}

func HandleGetAlertDetail(svc service.PestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		res, err := svc.GetReportDetail(r.Context(), id)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, "Laporan tidak ditemukan", nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "Detail Laporan", res)
	}
}

func HandleVerifyAlert(svc service.PestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := svc.VerifyReport(r.Context(), id); err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "Laporan berhasil diverifikasi", nil)
	}
}