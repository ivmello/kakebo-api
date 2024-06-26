package transactions

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ivmello/kakebo-go-api/internal/utils"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input CreateTransactionInput
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.service.CreateTransaction(ctx, input)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}

func (h *handler) ListAllUserTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	output, err := h.service.GetAllUserTransactions(ctx, userId)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}

func (h *handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	transactionId := r.PathValue("id")
	output, err := h.service.GetTransaction(ctx, userId, transactionId)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}

func (h *handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	transactionId := r.PathValue("id")
	err := h.service.DeleteTransaction(ctx, userId, transactionId)
	if err != nil {
		utils.JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, nil, http.StatusNoContent)
}
