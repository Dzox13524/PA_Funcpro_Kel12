package handle

import (
	"encoding/json"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/response"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/service"
)

func HandleCreateReservation(createSvc service.CreateReservationFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())

		var req domain.CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid Request Body", nil)
			return
		}

		res, err := createSvc(r.Context(), userID, req)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusCreated, "Reservasi berhasil dibuat", res)
	}
}

func HandleCreateOrder(createOrderSvc service.CreateOrderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())

		var req domain.CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid Request Body", nil)
			return
		}

		res, err := createOrderSvc(r.Context(), userID, req)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusCreated, "Order berhasil dibuat", res)
	}
}

func HandleGetMyReservations(getSvc service.GetUserTransactionsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())

		res, err := getSvc(r.Context(), userID)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusOK, "List Reservasi Saya", res)
	}
}

func HandleConfirmReservation(updateStatusSvc service.UpdateTransactionStatusFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if err := updateStatusSvc(r.Context(), id, domain.StatusConfirmed); err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "Reservasi dikonfirmasi", nil)
	}
}

func HandleCancelReservation(updateStatusSvc service.UpdateTransactionStatusFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if err := updateStatusSvc(r.Context(), id, domain.StatusCancelled); err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.WriteJSON(w, http.StatusOK, "Reservasi dibatalkan", nil)
	}
}

func HandleGetMyOrders(getSvc service.GetUserTransactionsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())

		res, err := getSvc(r.Context(), userID)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusOK, "List Pesanan Saya", res)
	}
}

func HandleGetTransactionDetail(getDetailSvc service.GetTransactionDetailFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		res, err := getDetailSvc(r.Context(), id)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, "Transaksi tidak ditemukan", nil)
			return
		}

		response.WriteJSON(w, http.StatusOK, "Detail Transaksi", res)
	}
}

func HandleUpdateOrderStatus(updateStatusSvc service.UpdateTransactionStatusFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var req domain.UpdateStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}

		if err := updateStatusSvc(r.Context(), id, req.Status); err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.WriteJSON(w, http.StatusOK, "Status pesanan diperbarui", nil)
	}
}