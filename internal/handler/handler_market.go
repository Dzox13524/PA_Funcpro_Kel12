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
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Unauthorized")
		}

		var req domain.CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid Request Body")
		}

		res, err := createSvc(r.Context(), userID, req)
		if err != nil {
			return 0, nil, err
		}

		return http.StatusCreated, res, nil
	})
}

func HandleCreateOrder(createOrderSvc service.CreateOrderFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Unauthorized")
		}

		var req domain.CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid Request Body")
		}

		res, err := createOrderSvc(r.Context(), userID, req)
		if err != nil {
			return 0, nil, err
		}

		return http.StatusCreated, res, nil
	})
}

func HandleGetMyReservations(getSvc service.GetUserTransactionsFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Unauthorized")
		}

		res, err := getSvc(r.Context(), userID)
		if err != nil {
			return 0, nil, err
		}

		return http.StatusOK, res, nil
	})
}

func HandleConfirmReservation(updateStatusSvc service.UpdateTransactionStatusFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")

		if err := updateStatusSvc(r.Context(), id, domain.StatusConfirmed); err != nil {
			return 0, nil, err
		}
		
		return http.StatusOK, nil, nil
	})
}

func HandleCancelReservation(updateStatusSvc service.UpdateTransactionStatusFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")

		if err := updateStatusSvc(r.Context(), id, domain.StatusCancelled); err != nil {
			return 0, nil, err
		}
		
		return http.StatusOK, nil, nil
	})
}

func HandleGetMyOrders(getSvc service.GetUserTransactionsFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		userID := middleware.GetUserIDFromContext(r.Context())
		if userID == "" {
			return 0, nil, response.NewAPIError(http.StatusUnauthorized, "Unauthorized")
		}

		res, err := getSvc(r.Context(), userID)
		if err != nil {
			return 0, nil, err
		}

		return http.StatusOK, res, nil
	})
}

func HandleGetTransactionDetail(getDetailSvc service.GetTransactionDetailFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")

		res, err := getDetailSvc(r.Context(), id)
		if err != nil {
			return 0, nil, response.NewAPIError(http.StatusNotFound, "Transaksi tidak ditemukan")
		}

		return http.StatusOK, res, nil
	})
}

func HandleUpdateOrderStatus(updateStatusSvc service.UpdateTransactionStatusFunc) http.HandlerFunc {
	return middleware.MakeHandler(func(r *http.Request) (int, any, error) {
		id := r.PathValue("id")

		var req domain.UpdateStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return 0, nil, response.NewAPIError(http.StatusBadRequest, "Invalid JSON")
		}

		if err := updateStatusSvc(r.Context(), id, req.Status); err != nil {
			return 0, nil, err
		}

		return http.StatusOK, nil, nil
	})
}