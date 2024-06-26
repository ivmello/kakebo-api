package transactions

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ivmello/kakebo-go-api/internal/core/transactions/dto"
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
	var input dto.CreateTransactionInput
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.service.CreateTransaction(ctx, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}

func (h *handler) ListAllUserTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	output, err := h.service.GetAllUserTransactions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}

func (h *handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transactionId := r.PathValue("id")
	output, err := h.service.GetTransaction(ctx, transactionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, output, http.StatusOK)
}
